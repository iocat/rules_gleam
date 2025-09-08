-module(hello_ffi).
-export([naive/1, tail_recursive/1]).

%% ====================================================================
%% Tail-Recursive Implementation
%% ====================================================================
-spec tail_recursive(non_neg_integer()) -> non_neg_integer().
tail_recursive(N) when is_integer(N), N >= 0 ->
    %% We start the helper with the 0th and 1st Fibonacci numbers as accumulators.
    tail_recursive_helper(N, 0, 1).



-spec tail_recursive_helper(non_neg_integer(), non_neg_integer(), non_neg_integer()) -> non_neg_integer().
tail_recursive_helper(0, A, _) -> A;
tail_recursive_helper(N, A, B) when N > 0 ->
    tail_recursive_helper(N - 1, B, A + B).
