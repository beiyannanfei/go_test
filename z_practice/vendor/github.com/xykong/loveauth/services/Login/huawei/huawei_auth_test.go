package huawei

import (
	"fmt"
	"testing"
)

func TestHuaweiAuth(t *testing.T) {
	result, err := AuthToken("1179011869901016", "1",
		"bVnFkJQRwXOKrpfRl9g1yni2EYz7/dIi/VOOI65VPjX2IKo+mm9TXf5knu4dmgIiIcxUDEzmgYhMfz9+LxQSRbmUkP+HhMUH+JOKgnBrSoNGcyDrooFgE3t0ETDtXu6kAqAUUbuM+oM2tzU9CaasbmDAkBaJ0HB/iyTfmke8IMvate7k3j3agzF7yxUsXnchsOmVRXwF8QwVBXD/xKK57kn3fn7Ei1SobSuvfH08iyweUNkVKUjqUlqkhg+GuR0Urz7czbYj6Zhi80+eKA4qxBKMb7d6qlaytAuzCrfuQWbpumatcFMo5FsMp2lJD9wYqsTAn7QZg1Bs6c3XlWKIjw==",
		"1547006509113")
	//"1350539672156")
	//"1546068044156")
	fmt.Printf("result = %#v\n", result)
	fmt.Printf("err = %v\n", err)
}
