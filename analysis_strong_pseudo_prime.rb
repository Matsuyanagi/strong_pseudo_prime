#!/usr/bin/env ruby
# -*- coding: utf-8 -*-
#-----------------------------------------------------------------------------
#	底と強擬素数の関係をパースして調べる
#	
#	
#	2017-07-20
#-----------------------------------------------------------------------------
require 'pp'
require 'fileutils'
require 'date'
require 'time'
require 'strscan'
require 'prime'

Encoding.default_external="utf-8"
#-----------------------------------------------------------------------------
#	
#-----------------------------------------------------------------------------
settings = {
	filename_input:	"strong_pseudo_prime.2_1024.3_10000000.txt",
	
}



#-----------------------------------------------------------------------------
#	
#-----------------------------------------------------------------------------
def main( settings )
	
	basenumber_and_pseudoprimes = []
	
	# ファイル読み込み
	filename_input = settings[ :filename_input ]
	
	# 底と、強擬素数の配列に読み込む
	File.foreach( filename_input ) do |line|
		if /^\s*(?<base_number>\d+) :[^\[]+\[(?<pseudo_numbers>[^\]]+)/ === line
			base_number = $~[:base_number].to_i
			pseudo_numbers = $~[:pseudo_numbers].split(',').map(&:strip).map(&:to_i)
			
			basenumber_and_pseudoprimes << [ base_number, pseudo_numbers ]
			
		end
		
	end
	
	# 配列を basenumber 順にソート
	basenumber_and_pseudoprimes = basenumber_and_pseudoprimes.sort_by{ |a| a[ 0 ] }
	# 強擬素数配列のカウントを追加
	basenumber_and_pseudoprimes.map!{ |a| a << a[1].count }
	# [ 0 ] basenumber
	# [ 1 ] pseudoprimes
	# [ 2 ] pseudoprimes.count
	#	例 : 
	#		[
	#			[ 2, [ 2047, 3277, 4033, 4681, 8321, ... ], ←の要素数 ],
	#			[ 3, [ 121, 703, 1891, 3281, 8401, ...   ], ←の要素数 ],
	#			...
	
	
	# pseudoprimes になり得た数値をすべて集めた
	pseudoprimes = basenumber_and_pseudoprimes.inject([]){ |p,a| p |= a[1]; p }.sort
	
	# 強擬素数の最小、最大、数
	puts "-- 最小、最大、重複をのぞいたユニークなカウント"
	puts pseudoprimes.min
	puts pseudoprimes.max
	puts pseudoprimes.count
	
	puts( "-- pseudoprimes" )
	puts( pseudoprimes.map(&:to_s).join(", ") )
	
	puts( "-- 数順" )
	basenumber_and_pseudoprimes.sort_by{ |a| [ a[ 2 ], a[ 0 ] ] }.each do |a|
		puts( "%4d : %s : %4d : %s" % [ a[ 0 ], a[ 0 ].prime? ? "*" : " " , a[ 2 ], a[ 0 ].prime_division ] )
	end
	
	puts( "-- 強擬素数の出現頻度 count > 10" )
	# pseudoprimes のヒストグラム
	pseudoprimes_histogram = Hash.new{ |h,k| h[k] = 0 }
	pseudoprimes_histogram_last1 = Hash.new{ |h,k| h[k] = 0 }
	basenumber_and_pseudoprimes.each do | base_number, pseudoprimes, pseudoprimes_count |
		pseudoprimes.each do |p|
			pseudoprimes_histogram[ p ] += 1
			pseudoprimes_histogram_last1[ p % 10 ] += 1
		end
	end
	# ハッシュを配列化して値ソート
	i = 0
	pseudoprimes_histogram = pseudoprimes_histogram.to_a.sort_by{ |k,v| [ -v, k, i+=1 ] }
	pseudoprimes_histogram.each do |pseudoprime, count|
		if count > 10
			puts( "%10d : %d" % [ pseudoprime, count ] )
		end
	end
	
	puts( "-- 強擬素数の末尾の出現頻度 強擬素数、末尾1桁の出現頻度" )
	# ハッシュを配列化してキーソート
	i = 0
	pseudoprimes_histogram_last1.to_a.sort_by{ |k,v| [ k, v, i+=1 ] }.each do |pseudoprime_last1, count|
		puts( "%d : %d" % [ pseudoprime_last1, count ] )
	end
	
	
	

end

main( settings )
