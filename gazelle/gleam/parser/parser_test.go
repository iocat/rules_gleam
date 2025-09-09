package parser

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// Mocking a testParse function, replace this with actual function from parser.go
func testParse(input string) (SourceFile, error) {
	result, err := ParseReader("test.gleam", strings.NewReader(input), Recover(true))
	if err != nil {
		return SourceFile{}, err
	}
	return result.(SourceFile), nil
}

func makeFunctionStmt(public bool, name string, params []Parameter, attrs []ExternalAttribute, ret Type) Node {
	if attrs == nil {
		attrs = []ExternalAttribute{}
	}
	if params == nil {
		params = []Parameter{}
	}
	return Statement(Function{
		Public:             public,
		Name:               name,
		Parameters:         params,
		ExternalAttributes: attrs,
		ReturnType:         ret,
	})
}

func makeImportStmt(module string, unqualifies ...string) Node {
	unqualifiedImports := []UnqualifiedImport{}
	for _, u := range unqualifies {
		unqualifiedImports = append(unqualifiedImports, UnqualifiedImport{Name: strings.TrimPrefix(u, "type:"), IsType: strings.HasPrefix(u, "type:")})
	}
	imp := Import{
		Module: module,
	}
	if len(unqualifiedImports) > 0 {
		imp.Unqualified = unqualifiedImports
	}
	return Statement(imp)
}

func makeParameters(params ...string) []Parameter {
	if len(params)%2 != 0 {
		panic("expected an even number of params")
	}

	var p []Parameter
	for i := 0; i < len(params); i += 2 {
		p = append(p, makeParameter(params[i], params[i+1]))
	}
	return p
}

func makeParameter(name, typ string) Parameter {
	return Parameter{Name: name, Type: string(typ)}
}


func TestParser(t *testing.T) {
	testCases := []struct {
		desc  string
		input string
		ast   SourceFile
	}{
		{
			desc: "main function",
			input: `
pub fn main() {
	echo "Hello, World!"
}`,
			ast: SourceFile{
				Statements: []Node{
					makeFunctionStmt(true, "main", nil, nil, nil),
				},
			},
		},
		{
			desc: "imports",
			input: `
import argv
import envoy
import filepath/*
import glam/doc.{type Document}
import gleam/bool
import gleam/dict.{type Dict}
import gleam/int
import gleam/io
import gleam/list
import gleam/option.{type Option, None, Some}

pub fn main() {
	echo "Hello, World!"
}`,
			ast: SourceFile{
				Statements: []Node{
					makeImportStmt("argv"),
					makeImportStmt("envoy"),
					makeImportStmt("filepath"),
					makeImportStmt("glam/doc", "type:Document"),
					makeImportStmt("gleam/bool"),
					makeImportStmt("gleam/dict", "type:Dict"),
					makeImportStmt("gleam/int"),
					makeImportStmt("gleam/io"),
					makeImportStmt("gleam/list"),
					makeImportStmt("gleam/option", "type:Option", "None", "Some"),
					makeFunctionStmt(true, "main", nil, nil, nil),
				},
			},
		},
		{
			desc: "external ffi",
			input: `
@external(erlang, "squirrel_ffi", "exit")
@external(javascript, "squirrel_ffi.mjs", "woah")
fn exit(n: Int) -> Nil

pub fn main() {
	exit(0)
}
`,
			ast: SourceFile{
				Statements: []Node{
					makeFunctionStmt(false, "exit", makeParameters("n", "Int"), []ExternalAttribute{
						{TargetLang: "erlang", Module: "squirrel_ffi", Function: "exit"},
						{TargetLang: "javascript", Module: "squirrel_ffi.mjs", Function: "woah"},
					}, string("Nil")),
					makeFunctionStmt(true, "main", nil, nil, nil),
				},
			},
		},
		{
			desc: "complex",
			input: `
import gleam/string
import gleam/result
import gleam/int

// @external is used to define a function that is implemented in another
// language, in this case, JavaScript. This is useful for interoperability.
// The first string is the target ("javascript"), the second is the module 
// ("./user_utils.mjs"), and the third is the function name ("log_external").
@external(javascript, "./user_utils.mjs", "log_external")
pub fn log_message(message: String) -> Nil

// A custom type representing possible validation errors for a user.
// The type is public, so it can be used by other modules, but its
// variants are defined here.
pub type ValidationError {
  InvalidName(reason: String)
  InvalidAge(reason: String)
}

// A public custom type for a User. This is a record type.
pub type User {
  User(name: String, age: Int)
}

pub fn create_user(name: String, age: Int) -> Result(User, ValidationError) {
  case is_valid_name(name) {
    Ok(valid_name) -> {
      case is_valid_age(age) {
        Ok(valid_age) -> Ok(User(name: valid_name, age: valid_age))
        Error(reason) -> Error(InvalidAge(reason))
      }
    }
    Error(reason) -> Error(InvalidName(reason))
  }
}

fn is_valid_name(name: String) -> Result(String, String) {
  if string.is_empty(name) {
    Error("Name cannot be empty.")
  } else if string.length(name) > 50 {
    Error("Name cannot be longer than 50 characters.")
  } else {
    Ok(name)
  }
}

// Another private function, this time for validating the age.
fn is_valid_age(age: Int) -> Result(Int, String) {
  if age <= 0 {
    Error("Age must be a positive number.")
  } else if age > 120 {
    Error("Age seems unlikely, please check.")
  } else {
    Ok(age)
  }
}

// A public helper function to get a greeting message for a user.
pub fn greeting(user: User) -> String {
  "Hello, " <> user.name <> "!"
}

// An example of how another module might use this one.
// We can't run this directly without a main project setup,
// but it shows the intended usage.
pub fn main() {
  // Let's try creating a valid user
  let valid_user_result = create_user("Alice", 30)

  case valid_user_result {
    Ok(user) -> {
      let message = greeting(user)
      log_message("Successfully created user: " <> message)
    }
    Error(error) -> {
      let error_message = case error {
        InvalidName(reason) -> "Invalid name: " <> reason
        InvalidAge(reason) -> "Invalid age: " <> reason
      }
      log_message("Failed to create user: " <> error_message)
    }
  }

  // Now, let's try an invalid one
  let invalid_user_result = create_user("", -5)
  case invalid_user_result {
    Ok(_) -> log_message("This should not have happened.")
    Error(error) -> {
      let error_message = case error {
        InvalidName(reason) -> "Invalid name: " <> reason
        InvalidAge(reason) -> "Invalid age: " <> reason
      }
      log_message("Correctly failed to create user: " <> error_message)
    }
  }
}
`,
			ast: SourceFile{Statements: []Node{
				makeImportStmt("gleam/string"),
				makeImportStmt("gleam/result"),
				makeImportStmt("gleam/int"),
				makeFunctionStmt(true, "log_message", []Parameter{
					makeParameter("message", "String"),
				}, []ExternalAttribute{
					{TargetLang: "javascript", Module: "./user_utils.mjs", Function: "log_external"},
				}, string("Nil")),
				makeFunctionStmt(true, "create_user",
					makeParameters("name", "String", "age", "Int"), nil, string("Result")),
				makeFunctionStmt(false, "is_valid_name",
					makeParameters("name", "String"), nil, string("Result")),
				makeFunctionStmt(false, "is_valid_age",
					makeParameters("age", "Int"), nil, string("Result")),
				makeFunctionStmt(true, "greeting",
					makeParameters("user", "User"), nil, string("String")),
				makeFunctionStmt(true, "main", nil, nil, nil),
			}},
		},
	}

	for _, tc := range testCases {
		ast, err := testParse(tc.input)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if diff := cmp.Diff(tc.ast, ast); diff != "" {
			t.Fatalf("desc: (%s)\n(-want, +got)=\n%s", tc.desc, diff)
		}
	}
}
