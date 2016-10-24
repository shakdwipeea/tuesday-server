package tuesday

import "gopkg.in/redis.v5"

// SetupDB for Identifier generation
// For details on discussions of identifier generation  refer https://github.com/shakdwipeea/tuesday-android/issues/14
//
// For the first 45k people the unique id wiil be
//	3 pairs of consecutive digits/alphabets
//	Eg AB ST 67
//	   34 GH 12
//	   ZA BC QR
//
// Here, we have to keep track of all the alphabets and numbers appearing in each position.
//
// So the db organization may look something like
// {
//	"first": {
//		"nextChar": "[A-Z]",
//		"nextNum": "[0-9]"
//	},
//	"second": {
//		"nextChar": "[A-Z]",
//		"nextNum": "[0-9]"
//	}
//	"third": {
//		"nextChar": "[A-Z]",
//		"nextNum": "[0-9]"
//	},
//	"numbersGenerated": []
// }
//
// their should be utility functions
//	to traverse the alphabets and the numbers
//	to find if a number has been used
//	randomly select if the current set is going to be alphabets or number
//
// make sure no collisions for 45k numbers, then we can remove the check and asynchronously notify user for any errors
// since we have established no chances for error.
func SetupDB(client *redis.Client) error {
	if err := client.HMSet("first", map[string]string{
		"nextChar": "a",
		"nextNum": "0",
	}).Err(); err != nil {
		return err
	}

	if err := client.HMSet("second", map[string]string{
		"nextChar": "s",
		"nextNum": "0",
	}).Err(); err != nil {
		return err
	}

	return client.HMSet("third", map[string]string{
		"nextChar": "z",
		"nextNum": "9",
	}).Err()
}