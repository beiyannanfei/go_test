package query

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/errors"
	"github.com/xykong/loveauth/settings"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/utils"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

func init() {
	handlers["/config"] = config
	getHandlers["/config"] = config

	go GlobalConfigSettingsWatcher()
}

// Binding from JSON
// swagger:parameters query_config_get
type DoQueryConfigReq struct {
	Token    string `form:"token" json:"token" binding:"required"`
	Version  string `form:"version" json:"version" binding:"required,min=5,max=40"`
	Branch   string `form:"branch" json:"branch" binding:"required"`
	Platform string `form:"platform" json:"platform"`
	Channel  string `form:"channel" json:"channel"`
}

//
// in: body
// swagger:parameters query_config
type DoQueryConfigReqBodyParams struct {
	//
	// swagger:allOf
	// in: body
	Request DoQueryConfigReq `form:"Request" json:"Request" binding:"required"`
}

// A DoQueryConfigRsp is an response message to client.
// swagger:response DoQueryConfigRsp
type DoQueryConfigRsp struct {
	// in: body
	Body struct {
		Code    int64  `json:"code"`
		Message string `json:"message"`

		// swagger:allOf
		// in: body
		Config GlobalConfigSettings `form:"config" json:"config" binding:"required"`
	}
}

type GlobalConfigSettings map[string]GlobalConfigSetting

// swagger:route GET /query/config query query_config_get
//
// Query config received from client.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Responses:
//       200: DoQueryConfigRsp

// swagger:route POST /query/config query query_config
//
// Query config received from client.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Responses:
//       200: DoQueryConfigRsp
func config(c *gin.Context) {

	var request DoQueryConfigReq

	// validation
	if err := c.Bind(&request); err != nil {

		utils.QuickReply(c, errors.ServerFail, "连接服务器异常")
		return
	}

	// validation version
	re := regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)$`)
	result := re.FindStringSubmatch(request.Version)
	if result == nil {
		utils.QuickReply(c, errors.ServerFail, "连接服务器异常")
		return
	}

	if request.Token == "" {

		utils.QuickReply(c, errors.ServerFail, "连接服务器异常")
		return
	}

	tokenRecord, err := storage.QueryAccessToken(request.Token)
	if err != nil {

		code := errors.Failed
		if ec, ok := err.(*errors.Type); ok {
			code = ec.Code
		}

		utils.QuickReply(c, code, "连接服务器异常")
		return
	}

	if tokenRecord.GlobalId == 0 {

		utils.QuickReply(c, errors.QueryAccessTokenFailed, "登录信息已失效")
		return
	}

	usedConfigs := matchConfig(request.Version, request.Platform, request.Channel, tokenRecord.GlobalId, request.Branch)
	if len(*usedConfigs) == 0 {

		logrus.WithFields(logrus.Fields{
			"version":  request.Version,
			"globalId": tokenRecord.GlobalId,
		}).Error("query_config can't matchConfig")

		utils.QuickReply(c, errors.ServerFail, "连接服务器异常")
		return
	}
	resp := DoQueryConfigRsp{}
	resp.Body.Code = int64(errors.Ok)
	resp.Body.Message = "Query config successfully!"
	resp.Body.Config = *usedConfigs
	c.JSON(http.StatusOK, resp.Body)
}

type GlobalConfigSetting struct {
	VersionRange   string        `json:"-"`
	ShowServerList bool          `json:"showServerList"`
	GVersion       string        `json:"gVersion"`
	Type           string        `json:"type"`
	BranchName     string        `json:"branchName"`
	LowVersion     string        `json:"lowVersion"`
	HighVersion    string        `json:"highVersion"`
	VersionList    []string      `json:"versionList"`
	Platform       []string      `json:"platform"`
	AllowChannel   []string      `json:"allowChannel"`
	RejectChannel  []string      `json:"rejectChannel"`
	Tags           []string      `json:"tags"`
	GConfigList    []interface{} `json:"gConfigList"`
}

//type GlobalConfig struct {
//	ConfigKey      string `json:"configKey"`
//	GlobalConfigVO struct {
//		CompleteBundle   int    `json:"completeBundle"`
//		ConfigURL        string `json:"configUrl"`
//		AccountURL       string `json:"accountUrl"`
//		BundleCdnURL     string `json:"bundleCdnUrl"`
//		BundleCdnURLRaw  string `json:"bundleCdnUrlRaw"`
//		GameServerURL    string `json:"gameServerUrl"`
//		HashNameFLag     int    `json:"hashNameFLag"`
//		ResVersion       int    `json:"resVersion"`
//		GmOpen           int    `json:"gmOpen"`
//		TestLoginOpen    int    `json:"testLoginOpen"`
//		AvgShowTestUI    int    `json:"avgShowTestUi"`
//		UsherOpen        int    `json:"usherOpen"`
//		SpeedScaleMax    int    `json:"speedScaleMax"`
//		ForbidGuide      int    `json:"forbidGuide"`
//		UpdateFlag       int    `json:"updateFlag"`
//		DownloadURL      string `json:"downloadUrl"`
//		Announcement     int    `json:"announcement"`
//		Prefix           string `json:"prefix"`
//		AnnouncementUrl  string `json:"announcementUrl"`
//		UrgentGConfigUrl string `json:"urgentGConfigUrl"`
//	} `json:"globalConfigVO"`
//}

var ConfigSettings []GlobalConfigSetting
var DefaultSetting GlobalConfigSettings
var BranchSetting []GlobalConfigSetting

func loadConfig(filename string) {

	b, err := ioutil.ReadFile(filename) // just pass the file name
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
			"path":  filename,
		}).Error("loadConfig load failed.")
	}

	var Settings map[string]GlobalConfigSetting
	err = json.Unmarshal(b, &Settings)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":               err,
			"globalConfigContent": string(b),
		}).Error("loadConfig Unmarshal failed.")

		return
	}

	ConfigSettings = make([]GlobalConfigSetting, 0)
	DefaultSetting = make(GlobalConfigSettings)
	BranchSetting = make([]GlobalConfigSetting, 0)

	var keys []string
	for k, _ := range Settings {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, versionRange := range keys {

		value := Settings[versionRange]

		if versionRange == "default" {

			DefaultSetting[versionRange] = value

			continue
		}

		if !strings.Contains(versionRange, "_") {

			if value.GVersion == "2.0" && value.Type == "2" {

				value.VersionRange = versionRange
				ConfigSettings = append(ConfigSettings, value)

				continue
			}
			//BranchSetting[versionRange] = value
			if value.GVersion != "2.0" {

				value.BranchName = versionRange
			}
			BranchSetting = append(BranchSetting, value)

			continue
		}

		value.VersionRange = versionRange
		ConfigSettings = append(ConfigSettings, value)
	}

	logrus.WithFields(logrus.Fields{
		"path":    filename,
		"content": ConfigSettings,
		"default": DefaultSetting,
		"branch":  BranchSetting,
	}).Info("loadConfig global config template loaded.")
}

func GlobalConfigSettingsWatcher() {

	reloadSeconds := settings.GetInt("loveauth", "global_config.reload_seconds")
	configTemplatePath := settings.GetString("loveauth", "global_config.path")

	loadConfig(configTemplatePath)

	if reloadSeconds == 0 {

		return
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	err = watcher.Add(configTemplatePath)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
			"path":  configTemplatePath,
		}).Error("GlobalConfigSettingsWatcher add watcher failed.")
		return
	}

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:

				logrus.WithFields(logrus.Fields{
					"path":  configTemplatePath,
					"event": event,
				}).Info("GlobalConfigSettingsWatcher watcher signal Events.")

				if event.Op&fsnotify.Write == fsnotify.Write {
					loadConfig(configTemplatePath)
				}

			case err := <-watcher.Errors:

				logrus.WithFields(logrus.Fields{
					"error": err,
					"path":  configTemplatePath,
				}).Error("GlobalConfigSettingsWatcher watcher signal Errors.")

			}
		}
	}()

	go func() {

		var fileMd5 string

		for {

			file, err := os.Open(configTemplatePath)
			if err != nil {

				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Error("GlobalConfigSettingsWatcher file open.")
			}

			md5Hash := md5.New()
			_, err = io.Copy(md5Hash, file)
			if err != nil {

				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Error("GlobalConfigSettingsWatcher md5 copy file.")
			}

			newFileMd5 := hex.EncodeToString(md5Hash.Sum(nil))

			if newFileMd5 != fileMd5 && fileMd5 != "" {

				loadConfig(configTemplatePath)
			}

			fileMd5 = newFileMd5

			time.Sleep(time.Second * time.Duration(reloadSeconds))
		}
	}()

	<-done
}

func toLongVersion(version string) string {

	re := regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)$`)
	result := re.FindStringSubmatch(version)
	if result == nil {
		return ""
	}

	for k, v := range result {
		if k == 0 {
			continue
		}

		i, _ := strconv.Atoi(v)

		result[k] = fmt.Sprintf("%9d", i)
	}

	return strings.Join(result[1:], ".")
}

func matchVersionRange(config GlobalConfigSetting, version string) bool {

	versionRange := config.VersionRange
	if config.GVersion == "2.0" {

		versionRange = fmt.Sprintf("%s_%s", config.LowVersion, config.HighVersion)
	}

	result := strings.Split(versionRange, "_")

	lowerVersion := toLongVersion(result[0])
	upperVersion := toLongVersion(result[1])
	longVersion := toLongVersion(version)

	if strings.Compare(longVersion, lowerVersion) == -1 {
		return false
	}

	if strings.Compare(longVersion, upperVersion) == 1 {
		return false
	}

	return true
}

func matchVersionList(versionList []string, version string) bool {

	for _, v := range versionList {

		if v == version {

			return true
		}
	}

	return false
}

func matchPlatformChannel(platfrom, channel string, config GlobalConfigSetting) bool {

	if config.GVersion != "2.0" {

		return true
	}

	if !inConfigPlatform(config.Platform, platfrom) {

		return false
	}

	if inAllowChannel(config.AllowChannel, channel) && !inRejectChannel(config.RejectChannel, channel) {

		return true
	}

	return false
}

func inConfigPlatform(platforms []string, platfrom string) bool {

	for _, p := range platforms {

		if p == platfrom {

			return true
		}
	}

	return false
}

func inAllowChannel(allowChannel []string, channel string) bool {

	if len(allowChannel) == 0 {

		return true
	}

	for _, c := range allowChannel {

		if c == channel {

			return true
		}
	}

	return false
}

func inRejectChannel(rejectChannel []string, channel string) bool {

	for _, c := range rejectChannel {

		if c == channel {

			return true
		}
	}

	return false
}

func getUserTag(globalId int64) string {

	var tags []string

	accountTag := storage.QueryAccountTag(globalId)
	if accountTag != nil {

		tags = append(tags, accountTag.Tags)
	}

	profile, _ := storage.QueryProfile(globalId)
	if profile != nil {

		//tags = append(tags, fmt.Sprintf("Reg_%d|Login_%d", profile.Auth.Extra.RegChannel, profile.Auth.Extra.LoginChannel))
		tags = append(tags, profile.Auth.Channel)
	}

	return strings.Join(tags, "|")
}

func matchAccountTag(globalId int64, configTags []string) bool {

	userTag := getUserTag(globalId)

	if len(userTag) == 0 {

		if len(configTags) == 0 {

			return true
		}

		return false
	}

	for _, tag := range configTags {

		if strings.Contains(userTag, tag) {

			return true
		}
	}

	return false
}

func matchConfig(version, platform, channel string, globalId int64, branch string) *GlobalConfigSettings {

	var usedConfigs = make(GlobalConfigSettings)

	var defaultVersionSettings = make(GlobalConfigSettings)

	for _, value := range ConfigSettings {

		if !matchPlatformChannel(platform, channel, value) {

			continue
		}

		if !matchVersionRange(value, version) && !matchVersionList(value.VersionList, version) {

			continue
		}

		if len(value.Tags) == 0 {

			defaultVersionSettings[value.VersionRange] = value

			continue
		}

		if matchAccountTag(globalId, value.Tags) {

			usedConfigs[value.VersionRange] = value

			return &usedConfigs
		}
	}

	//if len(defaultVersionSettings) == 0 && len(DefaultSetting) != 0 {
	//
	//	defaultVersionSettings = DefaultSetting
	//}

	if len(defaultVersionSettings) != 0 {

		return &defaultVersionSettings
	}

	for _, value := range BranchSetting {

		if value.BranchName == branch {

			if matchPlatformChannel(platform, channel, value) {

				defaultVersionSettings[value.BranchName] = value
				return &defaultVersionSettings
			}
		}
	}

	//if value, ok := BranchSetting[branch]; ok {
	//
	//	defaultVersionSettings[branch] = value
	//
	//	return &defaultVersionSettings
	//}

	return &DefaultSetting
}
