package main

import (
	"fmt"
	"os"
	_ "sort"
	"strings"
	_ "sync"
	"time"
)

const (
	BASE_LOWER_BOUND         = 2
	BASE_UPPER_BOUND         = 128
	PSEUDO_PRIME_LOWER_BOUND = 3
	//	PSEUDO_PRIME_UPPER_BOUND = 1000000
	PSEUDO_PRIME_UPPER_BOUND = 100000
)

type BaseNumberAndPseudoPrimeNumbers struct {
	BaseNumber         int64
	PseudoPrimeNumbers []int64
}

//	ソート用
/*
type BaseNumberAndPseudoPrimeNumbersSlice []BaseNumberAndPseudoPrimeNumbers
func (b BaseNumberAndPseudoPrimeNumbersSlice) Len() int {
	return len(b)
}
func (b BaseNumberAndPseudoPrimeNumbersSlice) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}
func (b BaseNumberAndPseudoPrimeNumbersSlice) Less(i, j int) bool {
	return b[i].BaseNumber < b[i].BaseNumber
}
*/
//	sort.Slice(b, func(i, j int) bool {
//	        return b[i].BaseNumber < b[j].BaseNumber
//	})

func main() {
	/*
		var base int64 = 2
		var exp int64 = 4
		var mod int64 = 13
		fmt.Printf("%d : %d : %d\n", base, exp, mod)
		ans := powmod(base, exp, mod)
		fmt.Printf("%d\n", ans)
	*/

	start := time.Now()

	strong_pseudoprimes := make(map[int64][]int64)

	for probable_prime := int64(PSEUDO_PRIME_LOWER_BOUND); probable_prime < PSEUDO_PRIME_UPPER_BOUND; probable_prime+=2 {

		calc_pseudoprimes(strong_pseudoprimes, probable_prime, BASE_LOWER_BOUND, BASE_UPPER_BOUND)

	}
	end := time.Now()

	//	結果出力
	for base_number, pseudoprimes := range strong_pseudoprimes {

		str := ""
		for _, v := range pseudoprimes {
			str += fmt.Sprintf("%d,", v)
		}
		str = strings.TrimRight(str, ",") //右端の","を取り除く

		fmt.Printf("%4d : %4d : %s\n", base_number, len(pseudoprimes), str)
	}

	fmt.Fprintf(os.Stderr, "%.2f sec\n", float64((end.Sub(start)).Nanoseconds())/1000000000.0)

}

func powmod(base int64, exp int64, mod int64) int64 {

	if base >= mod {
		base = base % mod
	}
	//	base が mod の倍数、もしくは等しい場合を除いておく
	if base == 0 {
		return 0
	}

	var answer int64 = 1
	for exp > 0 {
		if exp&1 == 1 {
			answer = (answer * base) % mod
		}
		base = (base * base) % mod
		exp >>= 1
	}
	return answer

}

func miller_rabin_primality_test(probable_prime int64, base int64) int {

	// 互いに素でないなら判定しない(片方がもう片方の倍数)
	if (base >= probable_prime && base%probable_prime == 0) || (base < probable_prime && probable_prime%base == 0) {
		return -1
	}

	// ミラーラビンテストで使用する d を求める
	// d : (probable_prime-1) = ( d*(2*r) ) を求める
	miller_rabin_d := probable_prime - 1
	for miller_rabin_d&1 == 0 {
		miller_rabin_d >>= 1
	}

	t := miller_rabin_d
	y := powmod(base, t, probable_prime)
	for t != probable_prime-1 && y != 1 && y != probable_prime-1 {
		y = (y * y) % probable_prime
		t <<= 1
	}
	if y != probable_prime-1 && t&1 == 0 {
		return 0
	}
	return 1
}

func calc_pseudoprimes(strong_pseudoprimes map[int64][]int64, probable_prime int64, lower_bound int64, upper_bound int64) {
	mr_test_results := make(map[int64]int, upper_bound-lower_bound+1)

	for base := lower_bound; base <= upper_bound; base++ {
		mr_test_results[base] = miller_rabin_primality_test(probable_prime, base)
	}

	answers_count := 0
	answers_sum := 0
	for _, test_result := range mr_test_results {
		if test_result >= 0 {
			answers_count++
			if test_result == 1 {
				answers_sum++
			}
		}
	}

	if answers_sum == 0 {
		//		return :composit
	} else if answers_sum == answers_count {
		//		return :prime
	} else {
		//	pseudoprime
		for base_number, test_result := range mr_test_results {
			if test_result == 1 {
				if _, ok := strong_pseudoprimes[base_number]; !ok {
					strong_pseudoprimes[base_number] = make([]int64, 0, 1000)
				}
				strong_pseudoprimes[base_number] = append(strong_pseudoprimes[base_number], probable_prime)
			}

		}
	}

}
