import mixed_bin_erl_internal/internal

@external(erlang, "test_ffi", "naive")
pub fn fibo_naive(a: Int) -> Int

pub fn main(){
    echo internal.fibo_tail_recursive(20)
    echo fibo_naive(10)
}