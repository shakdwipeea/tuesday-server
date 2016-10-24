package tuesday

import "testing"

func TestNextChar(t *testing.T) {
	chars := []string{
		"a","b","c","d","e","f","g","h","i","j","k","l","m",
		"n","o","p","q","r","s","t","u","v","w","x","y","z",
	}

	for i, char := range chars {
		checkWith := 0
		if i < len(chars) - 1 {
			checkWith = i + 1
		}

		if char, err := nextChar(char);
			err != nil || char != chars[checkWith] {
			t.Error("expected " + chars[checkWith] + " got " + char)
		}
	}
}

func TestNextInt(t *testing.T) {
	nums := []string{
		"0","1","2","3","4","5","6","7","8","9",
	}

	for i, num := range nums {
		checkWith := 0
		if i < len(nums) - 1 {
			checkWith = i + 1
		}
		if num, err := nextInt(num); err != nil || num != nums[checkWith] {
			t.Error("expected " + nums[checkWith])
		}
	}
}
