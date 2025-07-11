package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/term"
)

func main() {
	var temp string
	var slct int
	var msg string
	var key string
	var err error

	for {
		fmt.Println("选择[1:加密] [2:解密]: ")
		fmt.Scanln(&temp)
		slct, err = strconv.Atoi(temp)
		if err != nil {
			panic("输入错误")
		}
		if slct != 1 && slct != 2 {
			panic("输入有误")
		}

		fmt.Println("请输入明文/密文：")
		reader := bufio.NewReader(os.Stdin)
		msg, _ = reader.ReadString('\n')
		msg = strings.Trim(msg, "\n")

		fmt.Println("请输入密钥(不显示): ")
		byteKey := []byte{}
		for i := 0; i < 2; i++ {
			byteKey, err = term.ReadPassword(int(syscall.Stdin))
			if err != nil {
				panic("密钥输入错误")
			}
			if len(byteKey) > 0 {
				break
			} else if i == 1 {
				// 最后一次不输入密钥
				fmt.Println()
				panic("没有输入密钥...")
			}
			fmt.Print("请输入密钥: ")
		}
		key = string(byteKey)
		fmt.Printf("key: %v\n", key)
		if len([]rune(key)) > len([]rune(msg)) {
			panic("明文/密文长度太短...")
		}

		if slct == 1 {
			fmt.Println("加密结果：")
			fmt.Println(en_myBubbleSort(key, msg))
		} else {
			fmt.Println("解密结果：")
			fmt.Println(de_mydecryp(key, msg))
		}
		fmt.Println()
		fmt.Println()
	}
}

// 加密算法：对密钥排序的过程中，对明文同步进行排序从而实现顺序混淆
func en_myBubbleSort(key string, m string) string {
	newKey := keyFill(m, key)
	key_full := []rune(newKey)
	m2Rune := []rune(m)
	n := len(key_full)
	if n != len(m2Rune) {
		panic("err length...")
	}
	for i := n; i >= 0; i-- {
		for j := 0; j < i-1; j++ {
			if key_full[j] > key_full[j+1] {
				key_full[j], key_full[j+1] = key_full[j+1], key_full[j]
				m2Rune[j], m2Rune[j+1] = m2Rune[j+1], m2Rune[j]
			}
		}
	}
	return string(m2Rune)
}

// 解密算法：将密文按照初始密钥补充拆分为矩阵形式并重组
func de_mydecryp(key string, m_en string) string {
	n := len([]rune(m_en))
	m := "" //解密的结果
	newKey := keyFill(m_en, key)

	key2Rune := []rune(newKey)
	var key2Rune_1 []rune
	for _, item := range key2Rune {
		key2Rune_1 = append(key2Rune_1, item)
	}
	//对key2Rune_1排序
	sort.Slice(key2Rune_1, func(i, j int) bool {
		return key2Rune_1[i] < key2Rune_1[j]
	})

	m_en2Rune := []rune(m_en)
	hash := make(map[rune][]rune, 0)
	//拆分
	for i := 0; i < n; i++ {
		hash[key2Rune_1[i]] = append(hash[key2Rune_1[i]], m_en2Rune[i])
	}
	//重组
	for i := 0; i < n; i++ {
		m += string(hash[key2Rune[i]][0])         //将对应切片的第一个值作为解密结果
		hash[key2Rune[i]] = hash[key2Rune[i]][1:] //去掉对应切片的第一个值
	}
	return m
}

// 将密钥填充为与明文同样的长度
func keyFill(m string, key string) string {
	m2Rune := []rune(m)
	n := len(m2Rune)
	len := len(key)
	if n == len {
		return key
	}
	newKey := ""
	if n/len > 1 {
		for i := 0; i < n/len; i++ {
			newKey += key
		}
		key2Rune := []rune(key)
		for i := 0; i < n%len; i++ {
			newKey += string(key2Rune[i])
		}
	} else {
		newKey += key
		key2Rune := []rune(key)
		for i := 0; i < n%len; i++ {
			newKey += string(key2Rune[i])
		}
	}
	return newKey
}
