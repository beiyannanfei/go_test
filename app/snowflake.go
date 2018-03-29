package app

//https://github.com/bwmarrin/snowflake
import (
	""
	"fmt"
)

//生成一个id，类似uuid
func GeneraterId() {
	epoch := 1522292178485  //毫秒级时间戳
	snowflake.Epoch = epoch //Custom Epoch

	nodeId := 1 //1~1023
	node, err := snowflake.NewNode(nodeId)
	if err != nil {
		fmt.Printf("make snowflake node err: %v", err)
		return
	}
	id := node.Generate()

	// Print out the ID in a few different ways.
	fmt.Printf("Int64  ID: %d\n", id)
	fmt.Printf("String ID: %s\n", id)
	fmt.Printf("Base2  ID: %s\n", id.Base2())
	fmt.Printf("Base64 ID: %s\n", id.Base64())

	// Print out the ID's timestamp
	fmt.Printf("ID Time  : %d\n", id.Time())

	// Print out the ID's node number
	fmt.Printf("ID Node  : %d\n", id.Node())

	// Print out the ID's sequence number
	fmt.Printf("ID Step  : %d\n", id.Step())

	// Generate and print, all in one.
	fmt.Printf("ID       : %d\n", node.Generate().Int64())
}
