%% @doc This module demonstrates two implementations of the Fibonacci sequence.
%% 1. naive/1: A simple, but inefficient, recursive approach.
%% 2. tail_recursive/1: An efficient, tail-recursive approach suitable for large numbers.
-module(test_ffi).
-export([naive/1, tail_recursive/1]).

%% ====================================================================
%% Naive Recursive Implementation
%% ====================================================================

-spec naive(non_neg_integer()) -> non_neg_integer().
naive(0) -> 0;
naive(1) -> 1;
naive(N) when N > 1 ->
    naive(N - 1) + naive(N - 2).