#!/bin/bash
try() {
  expected="$1"
  input="$2"

  ./9gc "$input" > tmp.s
  gcc -o tmp tmp.s
  ./tmp
  actual="$?"

  if [ "$actual" = "$expected" ]; then
    echo "$input => $actual"
  else
    echo "$input => $expected expected, but got $actual"
    exit 1
  fi
}

try 0 0
try 42 42
try 41 " 12 + 34 - 5"
try 47 "5+6*7"
try 15 "5*(9-6)"
try 4 "(3+5)/2"
try 10 "-10+20"
try 10 "- -10"
try 10 "- - +10"
try 0 "0==1"
try 1 "42==42"
try 1 "0!=1"
try 0 "42!=42"

try 1 "0<1"
try 0 "1<1"
try 0 "2<1"
try 1 "0<=1"
try 1 "1<=1"
try 0 "2<=1"

try 1 "1>0"
try 0 "1>1"
try 0 "1>2"
try 1 "1>=0"
try 1 "1>=1"
try 0 "1>=2"

try 2 "a=2;"
try 10 "a=2;c=10"
try 99 "a=1;z=99"
try 32 "a=1;z=10;32"
try 57 "triple=3; nineteen=19; nineteen*triple"
try 5 "return 5
return 10;"
try 14 "
a = 3
b = 5 * 6 - 8
return a + b / 2
"

echo OK