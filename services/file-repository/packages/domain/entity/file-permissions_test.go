package entity

import (
	"testing"
)

// Test FilePermissionGroup constants
func TestFilePermissionGroupConstants(t *testing.T) {
	tests := []struct {
		permission FilePermissionGroup
		expected   uint16
	}{
		{DeleteFilePermission, 1 << 0},
		{UpdateFilePermission, 1 << 1},
		{ReadFilePermission, 1 << 2},
	}

	for _, test := range tests {
		if uint16(test.permission) != test.expected {
			t.Errorf("Expected %d, got %d", test.expected, test.permission)
		}
	}
}

// Test FilePermissionGroup String() method for all 8 combinations
func TestFilePermissionGroupString(t *testing.T) {
	tests := []struct {
		permission FilePermissionGroup
		expected   string
	}{
		// Single permissions
		{0, "---"},
		{ReadFilePermission, "r--"},
		{UpdateFilePermission, "-u-"},
		{DeleteFilePermission, "--d"},

		// Double permissions
		{ReadFilePermission | UpdateFilePermission, "ru-"},
		{ReadFilePermission | DeleteFilePermission, "r-d"},
		{UpdateFilePermission | DeleteFilePermission, "-ud"},

		// All permissions
		{ReadFilePermission | UpdateFilePermission | DeleteFilePermission, "rud"},
	}

	for _, test := range tests {
		result := test.permission.String()
		if result != test.expected {
			t.Errorf("For permission %b, expected '%s', got '%s'", test.permission, test.expected, result)
		}
	}
}

// Test ParseFilePermissionGroup for all valid combinations
func TestParseFilePermissionGroup_Valid(t *testing.T) {
	tests := []struct {
		input    string
		expected FilePermissionGroup
	}{
		{"---", 0},
		{"r--", ReadFilePermission},
		{"-u-", UpdateFilePermission},
		{"--d", DeleteFilePermission},
		{"ru-", ReadFilePermission | UpdateFilePermission},
		{"r-d", ReadFilePermission | DeleteFilePermission},
		{"-ud", UpdateFilePermission | DeleteFilePermission},
		{"rud", ReadFilePermission | UpdateFilePermission | DeleteFilePermission},
	}

	for _, test := range tests {
		result, err := ParseFilePermissionGroup(test.input)
		if err != nil {
			t.Errorf("Failed to parse '%s': %v", test.input, err)
			continue
		}
		if result != test.expected {
			t.Errorf("For input '%s', expected %b, got %b", test.input, test.expected, result)
		}
	}
}

// Test ParseFilePermissionGroup error cases
func TestParseFilePermissionGroup_Invalid(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{"Empty string", "", true},
		{"Too short", "ab", true},
		{"Too long", "abcd", true},
		{"Invalid char 1", "x--", true},
		{"Invalid char 2", "-x-", true},
		{"Invalid char 3", "--x", true},
		{"Multiple invalid", "xyz", true},
		{"Mixed case", "R--", true},
		{"Wrong order", "ur-", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := ParseFilePermissionGroup(test.input)
			if test.expectError && err == nil {
				t.Errorf("Expected error for input '%s', but got none. Result: %b", test.input, result)
			}
			if !test.expectError && err != nil {
				t.Errorf("Unexpected error for input '%s': %v", test.input, err)
			}
		})
	}
}

// Test FilePermissions Getter methods
func TestFilePermissionsGetters(t *testing.T) {
	owner := ReadFilePermission | UpdateFilePermission
	shared := ReadFilePermission | DeleteFilePermission
	other := UpdateFilePermission

	p := NewFilePermissions(owner, shared, other)

	if p.GetOwnerPermissions() != owner {
		t.Errorf("Owner permissions mismatch: expected %b, got %b", owner, p.GetOwnerPermissions())
	}
	if p.GetSharedPermissions() != shared {
		t.Errorf("Shared permissions mismatch: expected %b, got %b", shared, p.GetSharedPermissions())
	}
	if p.GetOtherPermissions() != other {
		t.Errorf("Other permissions mismatch: expected %b, got %b", other, p.GetOtherPermissions())
	}
}

// Test FilePermissions String() and parsing round-trip
func TestFilePermissionsRoundTrip(t *testing.T) {
	tests := []struct {
		name        string
		owner       FilePermissionGroup
		shared      FilePermissionGroup
		other       FilePermissionGroup
		expectedStr string
	}{
		{
			"All permissions all groups",
			ReadFilePermission | UpdateFilePermission | DeleteFilePermission,
			ReadFilePermission | UpdateFilePermission | DeleteFilePermission,
			ReadFilePermission | UpdateFilePermission | DeleteFilePermission,
			"rudrudrud",
		},
		{
			"No permissions",
			0, 0, 0,
			"---------",
		},
		{
			"Mixed permissions",
			ReadFilePermission | UpdateFilePermission,
			ReadFilePermission | DeleteFilePermission,
			UpdateFilePermission,
			"ru-r-d-u-",
		},
		{
			"Read only everywhere",
			ReadFilePermission,
			ReadFilePermission,
			ReadFilePermission,
			"r--r--r--",
		},
		{
			"Different combinations",
			ReadFilePermission,
			UpdateFilePermission,
			DeleteFilePermission,
			"r---u---d",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Test String() method
			p := NewFilePermissions(test.owner, test.shared, test.other)
			strResult := p.String()
			if strResult != test.expectedStr {
				t.Errorf("String() mismatch: expected '%s', got '%s'", test.expectedStr, strResult)
			}

			// Test parsing round-trip
			parsed, err := ParseFilePermissions(strResult)
			if err != nil {
				t.Errorf("Failed to parse '%s': %v", strResult, err)
				return
			}

			// Verify all permissions match
			if parsed.GetOwnerPermissions() != test.owner {
				t.Errorf("Owner round-trip mismatch: expected %b, got %b",
					test.owner, parsed.GetOwnerPermissions())
			}
			if parsed.GetSharedPermissions() != test.shared {
				t.Errorf("Shared round-trip mismatch: expected %b, got %b",
					test.shared, parsed.GetSharedPermissions())
			}
			if parsed.GetOtherPermissions() != test.other {
				t.Errorf("Other round-trip mismatch: expected %b, got %b",
					test.other, parsed.GetOtherPermissions())
			}

			// Test that string of parsed matches original string
			if parsed.String() != test.expectedStr {
				t.Errorf("Round-trip string mismatch: expected '%s', got '%s'",
					test.expectedStr, parsed.String())
			}
		})
	}
}

// Test ParseFilePermissions error cases
func TestParseFilePermissions_Invalid(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{"Empty string", "", true},
		{"Too short", "abc", true},
		{"Too long", "abcdefghij", true},
		{"Invalid length 7", "abcdefg", true},
		{"Invalid length 8", "abcdefgh", true},
		{"Invalid char in group1", "x--r--r--", true},
		{"Invalid char in group2", "r--x--r--", true},
		{"Invalid char in group3", "r--r--x--", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := ParseFilePermissions(test.input)
			if test.expectError && err == nil {
				t.Errorf("Expected error for input '%s', but got none. Result: %b", test.input, result)
			}
			if !test.expectError && err != nil {
				t.Errorf("Unexpected error for input '%s': %v", test.input, err)
			}
		})
	}
}

// Test FilePermissions IsEmpty method
func TestFilePermissionsIsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		owner    FilePermissionGroup
		shared   FilePermissionGroup
		other    FilePermissionGroup
		expected bool
	}{
		{"All empty", 0, 0, 0, true},
		{"Owner has permission", ReadFilePermission, 0, 0, false},
		{"Shared has permission", 0, UpdateFilePermission, 0, false},
		{"Other has permission", 0, 0, DeleteFilePermission, false},
		{"All have permissions", ReadFilePermission, UpdateFilePermission, DeleteFilePermission, false},
		{"Mixed empty", ReadFilePermission, 0, DeleteFilePermission, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := NewFilePermissions(test.owner, test.shared, test.other)
			result := p.IsEmpty()
			if result != test.expected {
				t.Errorf("IsEmpty() for (%b, %b, %b): expected %t, got %t",
					test.owner, test.shared, test.other, test.expected, result)
			}
		})
	}
}

// Test edge cases and bit manipulation
func TestFilePermissionsBitManipulation(t *testing.T) {
	t.Run("Maximum permissions", func(t *testing.T) {
		maxPerms := ReadFilePermission | UpdateFilePermission | DeleteFilePermission
		p := NewFilePermissions(maxPerms, maxPerms, maxPerms)

		// Should be able to get back what we put in
		if p.GetOwnerPermissions() != maxPerms {
			t.Errorf("Max owner perms not preserved")
		}
		if p.GetSharedPermissions() != maxPerms {
			t.Errorf("Max shared perms not preserved")
		}
		if p.GetOtherPermissions() != maxPerms {
			t.Errorf("Max other perms not preserved")
		}
	})

	t.Run("Single bit groups", func(t *testing.T) {
		// Test that each group only affects its own bits
		p1 := NewFilePermissions(ReadFilePermission, 0, 0)
		p2 := NewFilePermissions(0, UpdateFilePermission, 0)
		p3 := NewFilePermissions(0, 0, DeleteFilePermission)

		if p1.GetSharedPermissions() != 0 || p1.GetOtherPermissions() != 0 {
			t.Errorf("Owner-only permissions leaked to other groups")
		}
		if p2.GetOwnerPermissions() != 0 || p2.GetOtherPermissions() != 0 {
			t.Errorf("Shared-only permissions leaked to other groups")
		}
		if p3.GetOwnerPermissions() != 0 || p3.GetSharedPermissions() != 0 {
			t.Errorf("Other-only permissions leaked to other groups")
		}
	})
}

// Test constants
func TestConstants(t *testing.T) {
	if FilePermissionGroupSize != 3 {
		t.Errorf("FilePermissionGroupSize should be 3, got %d", FilePermissionGroupSize)
	}
	if FilePermissionGroupsAmount != 3 {
		t.Errorf("FilePermissionGroupsAmount should be 3, got %d", FilePermissionGroupsAmount)
	}

	// Test permission characters
	if ReadFilePermissionChar != 'r' {
		t.Errorf("ReadFilePermissionChar should be 'r', got %c", ReadFilePermissionChar)
	}
	if UpdateFilePermissionChar != 'u' {
		t.Errorf("UpdateFilePermissionChar should be 'u', got %c", UpdateFilePermissionChar)
	}
	if DeleteFilePermissionChar != 'd' {
		t.Errorf("DeleteFilePermissionChar should be 'd', got %c", DeleteFilePermissionChar)
	}
	if NoFilePermissionChar != '-' {
		t.Errorf("NoFilePermissionChar should be '-', got %c", NoFilePermissionChar)
	}
}
