package tss_sdk

import "C"
import (
	"github.com/sirupsen/logrus"
	"time"
	"unsafe"
	"github.com/xykong/loveauth/errors"
	"math/rand"
	"fmt"
	"github.com/xykong/loveauth/settings"
)

/*

#cgo LDFLAGS: -ldl
#include "tss_sdk_anti.h"
#include <string.h>
#include <stdio.h>
#include <stdlib.h>
#include <dlfcn.h>
typedef const TssSdkBusiInterf *(*TssSdkGetBusiInterf)(const TssSdkInitInfo* init_info);
#define TSS_SDK_GET_COMM_INTERF_SYM "tss_sdk_get_comm_interf"

typedef int (*TssSdkReleaseBusiInterf)();
#define TSS_SDK_RELEASE_BUSI_INTERF_SYM "tss_sdk_release_busi_interf"

static int g_is_sdk_loaded = 0;
void* g_handle = NULL;

#define LIB_NAME "libtss_sdk"
#define SUFFIX "so"


static int tss_sdk_load_linux(const char *shared_lib_file)
{
   fprintf(stdout, "tss_sdk_load_linux lib %s.\n", shared_lib_file);
   fflush(stdout);

   char *error = NULL;

   if (g_handle != NULL)
   {
       // 已经加载，则不再加载
       // Has been loaded, no longer load
		fprintf(stderr, "Has been loaded, no longer load %s failed, %s.\n", shared_lib_file, error);
		fflush(stderr);
       return 0;
   }

   g_handle = dlopen(shared_lib_file, RTLD_NOW);
   error = dlerror();
   if (error)
   {
		fprintf(stderr, "Open shared lib %s failed, %s.\n", shared_lib_file, error);
		fflush(stderr);
       return -1;
   }

   fprintf(stdout, "tss_sdk_load_linux lib %s load ok.\n", shared_lib_file);
   fflush(stdout);

   return 0;
}


void* tss_sdk_get_syml_linux(const char *syml_name)
{
   void *func = NULL;
   char *error = NULL;

	fprintf(stderr, "tss_sdk_get_syml_linux: %s.\n", syml_name);
	fflush(stdout);

   if (g_handle == NULL)
   {
		fprintf(stderr, "tss_sdk_get_syml_linux: %s failed, handler is not initialized.\n", syml_name);
		fflush(stdout);

       return NULL;
   }

   func = dlsym(g_handle, syml_name);
   error = dlerror();
   if (error)
   {
		fprintf(stderr, "dlsym failed, %s.\n", error);
		fflush(stdout);
       return NULL;
   }


	fprintf(stderr, "tss_sdk_get_syml_linux: %s ok.\n", syml_name);
	fflush(stdout);

   return func;
}


int tss_sdk_load_impl(const char *shared_lib_file)
{
   int ret = 0;
   if (g_is_sdk_loaded)
   {
       return 0;
   }

   #if defined(WIN32) || defined(_WIN64)
   ret = tss_sdk_load_win(shared_lib_file);
   #else
   ret = tss_sdk_load_linux(shared_lib_file);
   #endif
   if (ret == 0)
   {
       g_is_sdk_loaded = 1;
   }

   return ret;
}


void* tss_sdk_get_syml_impl(const char *syml_name)
{
   void *func = NULL;
   if (!g_is_sdk_loaded)
   {
       return NULL;
   }

#if defined(WIN32) || defined(_WIN64)
   func = tss_sdk_get_syml_win(syml_name);
#else
   func = tss_sdk_get_syml_linux(syml_name);
#endif

   return func;
}


// 获取sdk通用的接口
// Get the sdk common interface
static const TssSdkBusiInterf* tss_sdk_get_comm_interf(const TssSdkInitInfo* init_info)
{
	fprintf(stderr, "tss_sdk_get_comm_interf: %s.\n", init_info->tss_sdk_conf_);
	fflush(stdout);

	TssSdkGetBusiInterf func = NULL;
	func = (TssSdkGetBusiInterf)tss_sdk_get_syml_impl(TSS_SDK_GET_COMM_INTERF_SYM);
	if (func != NULL)
	{
		return func(init_info);
	}

	fprintf(stderr, "tss_sdk_get_comm_interf: %s ok.\n", init_info->tss_sdk_conf_);
	fflush(stdout);

	return NULL;
}


const TssSdkBusiInterf* tss_sdk_load(const char *shared_lib_dir, const TssSdkInitInfo *init_info)
{
   char shared_lib_file[1024] = {0};
   int ret = 0;

   if (shared_lib_dir == NULL || init_info == NULL)
   {
       return NULL;
   }

   snprintf(shared_lib_file,
            sizeof(shared_lib_file),
            "%s/%s.%s",
            shared_lib_dir,
            LIB_NAME,
            SUFFIX);

   ret = tss_sdk_load_impl(shared_lib_file);
   if (ret != 0)
   {
       // load library fail
       return NULL;
   }

   return tss_sdk_get_comm_interf(init_info);
}


int tss_sdk_unload()
{
   int rc = 0;
   TssSdkReleaseBusiInterf func = (TssSdkReleaseBusiInterf)tss_sdk_get_syml_impl(TSS_SDK_RELEASE_BUSI_INTERF_SYM);
   if (func != NULL)
   {
       // 调用接口释放函数
       // Call Interface release function
       func();
   }

   #if defined(WIN32) || defined(_WIN64)
   rc = FreeLibrary(g_handle);
   if (!rc)
   {
       return -2;
   }
   #else
   rc = dlclose(g_handle);
   if (rc != 0)
   {
       return -2;
   }
   #endif

   g_is_sdk_loaded = 0;
   g_handle = NULL;

   fprintf(stdout, "tss_sdk_unload success.\n");

   return 0;
}


// 获取某一特定业务的接口
// get the interface of a particular business
typedef const void *(*TssSdkGetInterf)(const void *init_data);
const void* tss_sdk_get_busi_interf(const char *syml_name, const void *data)
{
	fprintf(stderr, "tss_sdk_get_busi_interf: %s.\n", syml_name);
	fflush(stdout);

   TssSdkGetInterf f = (TssSdkGetInterf)tss_sdk_get_syml_impl(syml_name);
   if (f != NULL)
   {
		fprintf(stderr, "tss_sdk_get_busi_interf: %s ok.\n", syml_name);
		fflush(stdout);

       return f(data);
   }

	fprintf(stderr, "tss_sdk_get_busi_interf: %s failed.\n", syml_name);
	fflush(stdout);

   // 没找到对应的函数
   // Did not find the corresponding function
   return NULL;
}


extern TssSdkProcResult TssSdkSendDataToClientV3Callback(TssSdkAntiSendDataInfoV3 *anti_data);

// callback function of sending anti data，realized by game developer, called by sdk
TssSdkProcResult on_send_data_to_client(const TssSdkAntiSendDataInfoV3 *send_data_info) {

    return TSS_SDK_PROC_OK;
}


int anti_add_user(const void* in_anti_interf_,
 					const char* openid_, int plat_id_, int world_id_, int role_id_,
					const char* client_ip_, const char* client_ver_) {

	const TssSdkAntiInterfV3 *anti_interf_ = (const TssSdkAntiInterfV3 *)in_anti_interf_;
	if( anti_interf_ == NULL ) {
		return TSS_SDK_PROC_INVALID_ARG;
	}

    // user succeeded in logging in channel server, callthe add_user_ interface of anti-hacking service
    TssSdkAntiAddUserInfoV3 user_info;
    memset(&user_info, 0, sizeof(user_info));
    user_info.openid_.openid_ = (unsigned char *)openid_;
    user_info.openid_.openid_len_ = strlen(openid_);

    // platid, 0: IOS, 1: Android
    user_info.plat_id_ = plat_id_;
    user_info.world_id_ = world_id_;
    user_info.role_id_ = role_id_;
    //user_info.client_ip_ = (unsigned int)client_ip_;
    //user_info.client_ver_ = (unsigned int)client_ver_;
    user_info.client_ip_ = 0;
    user_info.client_ver_ = 0;

    //TssSdkUserExtData user_ext_data;
    //memcpy(user_ext_data.user_ext_data_, &uid, sizeof(uid));
    //user_ext_data.ext_data_len_ = sizeof(uid);
    //user_info.user_ext_data_ = &user_ext_data;

    return anti_interf_->add_user_(&user_info);
}


int anti_del_user(const void* in_anti_interf_, const char* openid_, int plat_id_, int world_id_, int role_id_) {

	const TssSdkAntiInterfV3 *anti_interf_ = (const TssSdkAntiInterfV3 *)in_anti_interf_;
	if( anti_interf_ == NULL ) {
		return TSS_SDK_PROC_INVALID_ARG;
	}

    // when user logged out the channel server, call the del_user_ interface of anti-hacking service
    TssSdkAntiDelUserInfoV3 user_info;
    memset(&user_info, 0, sizeof(user_info));
    user_info.openid_.openid_ = (unsigned char *)openid_;
    user_info.openid_.openid_len_ = strlen(openid_);
    user_info.openid_.openid_[user_info.openid_.openid_len_] = 0;
    // platid, 0: IOS, 1: Android
    user_info.plat_id_ = plat_id_;
    user_info.world_id_ = world_id_;
    user_info.role_id_ = role_id_;
    ////unsigned int uid = user->get_uid();
    //TssSdkUserExtData user_ext_data;
    //memcpy(user_ext_data.user_ext_data_, &uid, sizeof(uid));
    //user_ext_data.ext_data_len_ = sizeof(uid);
    //user_info.user_ext_data_ = &user_ext_data;

    return anti_interf_->del_user_(&user_info);
}

int on_recv_anti_data(const void* in_anti_interf_,
 					const char* openid_, int plat_id_, int world_id_, int role_id_,
					const char* anti_data_, unsigned int anti_data_len) {

    printf("GsBusiHandler::on_recv_anti_data\n");

	const TssSdkAntiInterfV3 *anti_interf_ = (const TssSdkAntiInterfV3 *)in_anti_interf_;
	if( anti_interf_ == NULL ) {
		return TSS_SDK_PROC_INVALID_ARG;
	}

    // call the recv data interface of sdk anti-hacking service to recv package
    TssSdkAntiRecvDataInfoV3 pkg_info;
    memset(&pkg_info, 0, sizeof(pkg_info));
    pkg_info.openid_.openid_ = (unsigned char *) openid_;
    pkg_info.openid_.openid_len_ = strlen(openid_);
    pkg_info.openid_.openid_[pkg_info.openid_.openid_len_] = 0;
    // platid, 0: IOS, 1: Android
    pkg_info.plat_id_ = plat_id_;
    pkg_info.world_id_ = world_id_;
    pkg_info.role_id_ = role_id_;
    //unsigned int uid = user->get_uid();
    //TssSdkUserExtData user_ext_data;
    //memcpy(user_ext_data.user_ext_data_, &uid, sizeof(uid));
    //user_ext_data.ext_data_len_ = sizeof(uid);
    //pkg_info.user_ext_data_ = &user_ext_data;
    pkg_info.anti_data_ = (const unsigned char *) anti_data_;
    pkg_info.anti_data_len_ = anti_data_len;

    return anti_interf_->recv_anti_data_(&pkg_info);
}

int busi_tick(const TssSdkBusiInterf* in_busi_interf_) {

	const TssSdkBusiInterf *busi_interf_ = (const TssSdkBusiInterf *)in_busi_interf_;
	if( busi_interf_ == NULL ) {
		return TSS_SDK_PROC_INVALID_ARG;
	}

	busi_interf_->proc_();

	return 0;
}

*/
import "C"

type TssSdk struct {
	handler    unsafe.Pointer
	antiInterf unsafe.Pointer
	busiInterf *C.TssSdkBusiInterf
}

func Start() {

	var tssSdk TssSdk

	tssSdk.Load()

	go func() {

		for {
			tssSdk.Tick()
			time.Sleep(time.Millisecond * 10)
		}

	}()

	go func() {

		for {
			time.Sleep(time.Second * 1)

			roleId := rand.Intn(10000)

			tssSdk.AddUser(fmt.Sprintf("user-%v", roleId), 1, 1, roleId, "127.0.0.1", "1.0.1")

			time.Sleep(time.Second * 1)

			tssSdk.OnRecvAntiData(fmt.Sprintf("user-%v", roleId), 1, 1, roleId, "test")
		}

	}()

	time.Sleep(time.Second * 5)

	tssSdk.Unload()
}

func (t *TssSdk) Load() error {

	//logrus.Info("tss_sdk load.")

	var tssSdkConf = settings.GetString("tencent", "tss.configs")
	var sharedLibDir = settings.GetString("tencent", "tss.shared_lib_dir")

	var initData C.TssSdkInitInfo
	initData.unique_instance_id_ = C.uint(time.Now().UnixNano())
	initData.tss_sdk_conf_ = C.CString(tssSdkConf)
	defer C.free(unsafe.Pointer(initData.tss_sdk_conf_))

	var cSharedLibDir = C.CString(sharedLibDir)
	defer C.free(unsafe.Pointer(cSharedLibDir))

	t.busiInterf = C.tss_sdk_load((*C.TSS_TCHAR)(cSharedLibDir), &initData)
	if t.busiInterf == nil {
		logrus.WithFields(logrus.Fields{
			"shared_lib_dir": sharedLibDir,
			"busiInterf":     t.busiInterf,
		}).Error("tss tss_sdk_load failed.")
		return errors.New("tss tss_sdk_load failed")
	}

	logrus.WithFields(logrus.Fields{
		"shared_lib_dir": sharedLibDir,
		"busiInterf":     t.busiInterf,
	}).Info("tss tss_sdk_load success.")

	// get anti hacking service interface handle anti_interf_
	var initInfo C.TssSdkAntiInitInfoV3
	initInfo.send_data_to_client_ = C.TssSdkSendDataToClientV3(C.TssSdkSendDataToClientV3Callback)

	var symlName = "tss_sdk_get_anti_interf_v3"
	var cSymlName = C.CString(symlName)
	defer C.free(unsafe.Pointer(cSymlName))

	t.antiInterf = C.tss_sdk_get_busi_interf(cSymlName, unsafe.Pointer(&initInfo))
	if t.antiInterf == nil {
		logrus.WithFields(logrus.Fields{
			"shared_lib_dir": sharedLibDir,
			"syml_name":      symlName,
			"anti_interf_":   t.antiInterf,
		}).Error("tss tss_sdk_get_busi_interf failed.")
		return errors.New("tss tss_sdk_get_busi_interf failed")
	}

	logrus.WithFields(logrus.Fields{
		"shared_lib_dir": sharedLibDir,
		"anti_interf_":   t.antiInterf,
	}).Info("tss tss_sdk_get_busi_interf success.")

	return nil
}

func (t *TssSdk) Unload() {

	logrus.Info("tss_sdk unload.")

	C.tss_sdk_unload()
}

func (t *TssSdk) AddUser(openId string, platId int, worldId int, roleId int, clientIP string, clientVer string) error {

	if t.antiInterf == nil {
		return errors.New("tss is not initialized")
	}

	var cOpenId = C.CString(openId)
	defer C.free(unsafe.Pointer(cOpenId))
	var cClientIP = C.CString(clientIP)
	defer C.free(unsafe.Pointer(cClientIP))
	var cClientVer = C.CString(clientVer)
	defer C.free(unsafe.Pointer(cClientVer))

	ret := C.anti_add_user(t.antiInterf, cOpenId, C.int(platId), C.int(worldId), C.int(roleId), cClientIP, cClientVer)

	if ret != 0 {
		return errors.New("tss add user failed")
	}

	logrus.Info("anti_add_user: ", ret)

	return nil
}

func (t *TssSdk) DelUser(openId string, platId int, worldId int, roleId int) error {

	if t.antiInterf == nil {
		return errors.New("tss is not initialized")
	}

	var cOpenId = C.CString(openId)
	defer C.free(unsafe.Pointer(cOpenId))

	ret := C.anti_del_user(t.antiInterf, cOpenId, C.int(platId), C.int(worldId), C.int(roleId))

	logrus.Info("anti_del_user: ", ret)

	return nil
}

func (t *TssSdk) OnRecvAntiData(openId string, platId int, worldId int, roleId int, antiData string) error {

	if t.antiInterf == nil {
		return errors.New("tss is not initialized")
	}

	var cOpenId = C.CString(openId)
	defer C.free(unsafe.Pointer(cOpenId))
	var cAntiData = C.CString(antiData)
	defer C.free(unsafe.Pointer(cAntiData))

	ret := C.on_recv_anti_data(t.antiInterf, cOpenId, C.int(platId), C.int(worldId), C.int(roleId), cAntiData, C.uint(len(antiData)))

	if ret != 0 {
		return errors.New("tss onRecvAntiData failed")
	}

	logrus.Info("on_recv_anti_data: ", ret)

	return nil
}

func (t *TssSdk) Tick() {

	//logrus.Info("tss_sdk busiTick.")

	C.busi_tick(t.busiInterf)
}
