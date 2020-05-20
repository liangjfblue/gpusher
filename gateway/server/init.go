/**
 *
 * @author liangjf
 * @create on 2020/5/20
 * @version 1.0
 */
package server

type IServer interface {
	Init() error
	Run()
}
