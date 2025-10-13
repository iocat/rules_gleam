package parser

import (
	"regexp"
	"strconv"
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

func makeImportStmt(module string, unqualifiesOrTarget ...string) Node {
	target := ""
	alias := ""
	unqualifiedImports := []UnqualifiedImport{}
	for _, u := range unqualifiesOrTarget {
		if strings.HasPrefix(u, "target:") {
			target = strings.TrimPrefix(u, "target:")
		} else if strings.HasPrefix(u, "alias:") {
			alias = strings.TrimPrefix(u, "alias:")
		} else {
			unqualifiedImports = append(unqualifiedImports, UnqualifiedImport{Name: strings.TrimPrefix(u, "type:"), IsType: strings.HasPrefix(u, "type:")})
		}
	}
	imp := Import{
		Module: module,
	}
	if len(target) > 0 {
		imp.Target = &TargetAttribute{
			TargetLang: target,
		}
	}
	if len(alias) > 0 {
		imp.Alias = alias
	}
	if len(unqualifiedImports) > 0 {
		imp.Unqualified = unqualifiedImports
	}
	return Statement(imp)
}

func makeParameters(params ...interface{}) []Parameter {
	if len(params)%2 != 0 {
		panic("expected an even number of params")
	}

	var p []Parameter
	for i := 0; i < len(params); i += 2 {
		p = append(p, makeParameter(params[i].(string), params[i+1]))
	}
	return p
}

func makeParameter(name string, typ interface{}) Parameter {
	var label string = ""
	if strings.Contains(name, " ") {
		label = name[:strings.Index(name, " ")]
		name = name[strings.Index(name, " ")+1:]
	}
	if typ == nil {
		return Parameter{Label: label, Name: name}
	}
	return Parameter{Label: label, Name: name, Type: string(typ.(string))}
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
@target(javascript)
import houdini/internal/escape


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
					makeImportStmt("houdini/internal/escape", "target:javascript"),
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
			desc: "with lambda",
			input: `
pub fn main(test: fn(Int) -> Int) {
	echo "Hello, World!"
}

/// let result = decode.run(data, decoder)
/// assert result == Ok(SignUp(name: "Lucy", email: "lucy@example.com"))
///
///
pub fn subfield(
  field_path: List(name),
  field_decoder: Decoder(t),
  next: fn(t) -> Decoder(final),
) -> Decoder(final) {
	Decoder(function: fn(data) {
    let #(out, errors1) =
      index(field_path, [], field_decoder.function, data, fn(data, position) {
        let #(default, _) = field_decoder.function(data)
        #(default, [DecodeError("Field", "Nothing", [])])
        |> push_path(list.reverse(position))
      })
    let #(out, errors2) = next(out).function(data)
    #(out, list.append(errors1, errors2))
  })
}

fn fold_dict(
  acc: #(Dict(k, v), dynamic.List(dynamic.DecodeError)),
  key: dynamic.Dynamic,
  value: fn(),
  key_decoder: fn(Dynamic) -> #(k, List(DecodeError)),
  value_decoder: fn(Dynamic) -> #(v, List(DecodeError)),
) -> #(Dict(k, v), List(DecodeError)) {}

pub fn high_order_function(effect: fn(fn(msg) -> Nil) -> Nil) -> Effect(msg) {
  Effect(..empty, synchronous: [task])
}

`,
			ast: SourceFile{
				Statements: []Node{
					makeFunctionStmt(true, "main", makeParameters("test", nil), nil, nil),
					makeFunctionStmt(true, "subfield", makeParameters("field_path", "List", "field_decoder", "Decoder", "next", nil), nil, string("Decoder")),
					makeFunctionStmt(false, "fold_dict", makeParameters("acc", nil, "key", "Dynamic", "value", nil, "key_decoder", nil, "value_decoder", nil), nil, nil),
					makeFunctionStmt(true, "high_order_function", makeParameters("effect", nil), nil, string("Effect")),
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

pub fn one_of(
  first: Decoder(a),
  or alternatives: List(Decoder(a)),
) -> Decoder(a) {
  Decoder(function: fn(dynamic_data) {
    let #(_, errors) as layer = first.function(dynamic_data)
    case errors {
      [] -> layer
      [_, ..] -> run_decoders(dynamic_data, layer, alternatives)
    }
  })
}

pub fn from_list(list: List(#(k, v))) -> Dict(k, v) {
  from_list_loop(list, new())
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
				makeFunctionStmt(true, "one_of",
					makeParameters("first", "Decoder", "or alternatives", "List"), nil, string("Decoder")),
				makeFunctionStmt(true, "from_list", makeParameters("list", "List"), nil, string("Dict")),
			}},
		},
		{
			desc: "houdini",
			input: `@target(javascript)
import houdini/internal/escape_js as escape

@target(erlang)
import houdini/internal/escape_erl as escape

/// Escapes a string to be safely used inside an HTML document by escaping
/// the following characters:
///
///` +
				"///   - `<` becomes `&lt;`" +
				"///   - `>` becomes `&gt;`" +
				"///   - `&` becomes `&amp;`" +
				"///   - `\"` becomes `&quot;`" +
				"///   - `'` becomes `&#39;`." +
				`///
/// ## Examples
///
/// assert escape("wibble & wobble") == "wibble &amp; wobble"
/// assert escape("wibble > wobble") == "wibble &gt; wobble"
`, ast: SourceFile{
				Statements: []Node{
					makeImportStmt("houdini/internal/escape_js", "alias:escape", "target:javascript"),
					makeImportStmt("houdini/internal/escape_erl", "alias:escape", "target:erlang"),
				},
			},
		},
		{
			desc: "houdini",
			input: `
/// An alernative ok type used within Erlang/OTP
pub type Result2(data1, data2, error) {
  Ok(data1, data2)
  Error(error)
}

@target(erlang)
type Token =
  List(Nil)
`, ast: SourceFile{
				Statements: []Node{},
			},
		},
	}

	for _, tc := range testCases {
		ast, err := testParse(tc.input)
		if err != nil {
			matches := regexp.MustCompile(`^.*test\.gleam:(\d+):`).FindStringSubmatch(err.Error())
			if len(matches) == 0 {
				t.Fatalf("failed to get parse line: maybe a different error err: %v", err)
			}
			line, _ := strconv.Atoi(matches[1])

			t.Fatalf("expected no error, got %v, test line:\n\n\t%v", err, strings.Split(tc.input, "\n")[line-1])
		}
		if diff := cmp.Diff(tc.ast, ast); diff != "" {
			t.Fatalf("desc: (%s)\n(-want, +got)=\n%s", tc.desc, diff)
		}
	}
}
