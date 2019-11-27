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

try 0 "func main() {return 0}"
try 42 "func main() {return 42}"
try 41 "func main() {return 12 + 34 - 5}"
try 47 "func main() {return 5+6*7}"
try 15 "func main() {return 5*(9-6)}"
try 4 "func main() {return (3+5)/2}"
try 10 "func main() {return -10+20}"
try 10 "func main() {return - -10}"
try 10 "func main() {return - - +10}"
try 0 "func main() {return 0==1}"
try 1 "func main() {return 42==42}"
try 1 "func main() {return 0!=1}"
try 0 "func main() {return 42!=42}"

try 1 "func main() {return 0<1}"
try 0 "func main() {return 1<1}"
try 0 "func main() {return 2<1}"
try 1 "func main() {return 0<=1}"
try 1 "func main() {return 1<=1}"
try 0 "func main() {return 2<=1}"

try 1 "func main() {return 1>0}"
try 0 "func main() {return 1>1}"
try 0 "func main() {return 1>2}"
try 1 "func main() {return 1>=0}"
try 1 "func main() {return 1>=1}"
try 0 "func main() {return 1>=2}"

try 2 "func main() {return a=2;}"
try 10 "func main() {a=2;return c=10}"
try 99 "func main() {a=1;return z=99}"
try 32 "func main() {a=1;z=10; return 32}"
try 57 "func main() {triple=3; nineteen=19;return nineteen*triple}"
try 5 "
func main() {
  return 5
  return 10;
}"
try 14 "
func main() {
  a = 3
  b = 5 * 6 - 8
  return a + b / 2
}
"
try 1 "
func main() {
  a = 1
  b = 4
  if a == 1 {
    return a
  }
  return b
}
"
try 4 "
func main() { 
  a = 1
  b = 4
  if a != 1 {
    return a
  }
  return b
}
"
try 10 "
func main() {
  if 1>2 {
    return 0
  } else {
    return 10
  }
}
"
try 10 "
func main() {
  if 1>2 {
    return 0
  } else if 1==0 {
    return -1
  } else {
    return 10
  }
}
"
try 10 "
func main() {
  j = 0
  for i = 0; i < 10; i++ {
    j++
  }
  return j
}
"
try 1 "
func main() {
  i = 10
  for i > 1 {
    i--
  }
  return i
}
"
try 3 "
func main() {
  x=1;y=2
  if z:=x+y; z==3 {
    return z
  }
  return -1
}
"
try 7 "
func main() {
  n := 0
  for i := 0; i < 5; i++ {
    n++
    if i == 1 { n++ }
    if i == 3 { n++ }
  }
  return n
}
"
try 3 "
func main() {
  return add(1, 2)
}
func add(a, b) {
  return a + b
}
"
try 89 '
func main() { 
  return fib(10)
}
func fib(x) { 
  if x<=1 {return 1}
  return fib(x-1) + fib(x-2)
}
'
try 3 '
func main() {
  x = 3
  y = &x
  return *y
}
'

echo OK