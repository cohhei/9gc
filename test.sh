#!/bin/bash
try() {
  expected="$1"
  input="$2"

  ./9gc "$input" > tmp.s
  gcc -static -o tmp tmp.s
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

try 2 "func main() {var a int;a=2;return a;}"
try 10 "func main() {var a int;var c int;a=2;c=10return c}"
try 99 "func main() {var a int;var z int;a=1;z=99;return z}"
try 32 "func main() {var a int;var z int;a=1;z=10; return 32}"
try 57 "func main() {var triple int;var nineteen int;triple=3; nineteen=19;return nineteen*triple}"
try 5 "
func main() {
  return 5
  return 10;
}"
try 14 "
func main() {
  var a int
  var b int
  a = 3
  b = 5 * 6 - 8
  return a + b / 2
}
"
try 1 "
func main() {
  a := 1
  b := 4
  if a == 1 {
    return a
  }
  return b
}
"
try 4 "
func main() { 
  a := 1
  b := 4
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
  j := 0
  for i := 0; i < 10; i++ {
    j++
  }
  return j
}
"
try 1 "
func main() {
  i := 10
  for i > 1 {
    i--
  }
  return i
}
"
try 3 "
func main() {
  x:=1;y:=2
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
func add(a int, b int) int {
  return a + b
}
"
try 89 '
func main() { 
  return fib(10)
}
func fib(x int) int { 
  if x<=1 {return 1}
  return fib(x-1) + fib(x-2)
}
'
try 3 '
func main() {
  x := 3
  y := &x
  return *y
}
'
try 3 '
func main() {
  var x int
  var y *int
  y = &x;
  *y = 3
  return x
}
'
try 0 '
func main() {
  var a [2][3]int
  return 0
}
'
try 15 'func main() { var x [2]int; x[0]=3; x[1]=5; return x[0] * x[1]; }'
try 2 'func main() { var x [2][3]int; x[1][2]=2; return x[1][2]; }'
try 1 'func main() { var x [2][3]int; x[0][0]=1; y:=x; return y[0][0]; }'
try 0 'var x int; func main() { return x }'
try 3 'var x int; func main() { x=3; return x }'
try 0 'var x [4]int; func main() { x[0]=0; x[1]=1; x[2]=2; x[3]=3; return x[0] }'
try 1 'var x [4]int; func main() { x[0]=0; x[1]=1; x[2]=2; x[3]=3; return x[1] }'
try 2 'var x [4]int; func main() { x[0]=0; x[1]=1; x[2]=2; x[3]=3; return x[2] }'
try 3 'var x [4]int; func main() { x[0]=0; x[1]=1; x[2]=2; x[3]=3; return x[3] }'
try 1 'func main() { var b byte; b = 1; return b }'
try 97 'func main() { return "abc"[0] }'
try 98 'func main() { return "abc"[1] }'
try 99 'func main() { a := "abc"; return a[2] }'
try 0 'func main() { return "abc"[3] }'

echo OK