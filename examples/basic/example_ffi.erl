%% @doc This module demonstrates two implementations of the Fibonacci sequence.
%% 1. naive/1: A simple, but inefficient, recursive approach.
%% 2. tail_recursive/1: An efficient, tail-recursive approach suitable for large numbers.
-module(example_ffi).
-export([naive/1, tail_recursive/1]).

%% ====================================================================
%% Naive Recursive Implementation
%% ====================================================================

-spec naive(non_neg_integer()) -> non_neg_integer().
naive(0) -> 0;
naive(1) -> 1;
naive(N) when N > 1 ->
    naive(N - 1) + naive(N - 2).


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
