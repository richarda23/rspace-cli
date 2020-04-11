package cmd
import (
"os"
"encoding/json"
"fmt"
//"errors"
)
func exitWithStdErrMsg (message string) {
	fmt.Fprintln(os.Stderr, message)
	os.Exit(1)
}
func exitWithErr (err error) {
	exitWithStdErrMsg(err.Error()) 
}
func prettyMarshal(anything interface{}) string {
        bytes, _ := json.MarshalIndent(anything, "", "\t")
        return string(bytes)
}
