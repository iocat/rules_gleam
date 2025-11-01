#!/usr/bin/env -S deno
import { main } from "[[main_js_path]]";

function exit(code) {
  if (globalThis.Deno) {
    Deno.exit(code);
  } else {
    process.exit(code);
  }
}

main();
