package require

import (
	"sync"

	"github.com/baker-yuan/go-gateway/eosc"
)

var (
	_ eosc.IRequires = (*Manager)(nil)
)

type Manager struct {
	locker    sync.Mutex                     //
	workerIds eosc.Untyped[string, []string] // 存储依赖关系 key=ID value=该ID依赖的所有ID列表
	requireBy eosc.Untyped[string, []string] // 存储各id的依赖项 key=依赖ID value=依赖该ID的所有ID列表
}

func NewRequireManager() eosc.IRequires {
	return &Manager{
		locker:    sync.Mutex{},
		requireBy: eosc.BuildUntyped[string, []string](),
		workerIds: eosc.BuildUntyped[string, []string](),
	}
}

func (w *Manager) Set(id string, requiresIds []string) {
	// a 依赖于 b c d
	// b c d 被 a依赖
	// a -> b c d
	// workerIds a => [b c d]
	// requireBy
	// b => a
	// c => a
	// d => a

	w.locker.Lock()
	defer w.locker.Unlock()

	// 删除原有的依赖关系
	w.del(id)

	if len(requiresIds) > 0 {
		for _, rid := range requiresIds {
			d, has := w.requireBy.Get(rid)
			if !has {
				w.requireBy.Set(rid, []string{id})
			} else {
				w.requireBy.Set(rid, append(d, id))
			}
		}
		w.workerIds.Set(id, requiresIds)
	}
}

func (w *Manager) RequireBy(requireId string) []string {
	// 获取依赖requireID的所有id列表
	w.locker.Lock()
	ids, has := w.requireBy.Get(requireId)
	w.locker.Unlock()
	if has {
		return ids
	}
	return nil
}

func (w *Manager) Del(id string) {
	w.locker.Lock()
	w.del(id)
	w.locker.Unlock()
}

func (w *Manager) del(id string) {
	if rs, has := w.workerIds.Del(id); has {
		for _, rid := range rs {
			w.removeBy(id, rid)
		}
	}
}

func (w *Manager) removeBy(id string, requireId string) {
	if rs, has := w.requireBy.Get(requireId); has {
		for i, rid := range rs {
			if rid == id {
				rs = append(rs[:i], rs[i+1:]...)
				break
			}
		}
		if len(rs) == 0 {
			w.requireBy.Del(requireId)
		} else {
			w.requireBy.Set(requireId, rs)
		}
	}
}

func (w *Manager) Requires(id string) []string {
	// 根据id获取所有依赖的id列表
	w.locker.Lock()
	ids, has := w.workerIds.Get(id)
	w.locker.Unlock()
	if has {
		return ids
	}
	return nil
}

func (w *Manager) RequireByCount(requireId string) int {
	w.locker.Lock()
	rs, has := w.requireBy.Get(requireId)
	w.locker.Unlock()
	if has {
		return len(rs)
	}
	return 0
}
