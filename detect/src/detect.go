package detect

type Detect interface {
	Run(target *DTarget) error //执行探测

	Publish(result *DResult) //若探测到相关的应用，则用此发布探测结果
}
