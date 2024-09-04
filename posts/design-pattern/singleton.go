//=============================================================================
/*!
 *  @file       singleton.go
 *  @brief      Singleton 单例模式
 *  @author     kydenlu
 *  @date       2024.09
 *  @note
 */
//=============================================================================

package designpattern

import "sync"

var once sync.Once

type Singleton struct {
	str string
}

var instance *Singleton

func GetInstance() *Singleton {
	if instance != nil {
		return instance
	}

	once.Do(func() {
		instance = &Singleton{}
	})
	return instance
}
