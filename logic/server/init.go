/**
 *
 * @author liangjf
 * @create on 2020/6/3
 * @version 1.0
 */
package server

type IServer interface {
	Init() error
	Run()
	Stop()
}
