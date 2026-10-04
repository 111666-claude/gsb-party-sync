// Package partysync 是队伍状态版本向量：逐成员序号、幂等应用、乱序暂存补洞、差异检测与修复。
// 缺陷：应用不幂等、旧序号会覆盖新状态、跳号直接接受、差异检测拿全局最大序号比较、
// 落后方向也报全部成员、每次应用都扫全部历史。
package partysync

import "sort"

// Vector 是成员名到已确认最新序号的版本向量。
type Vector map[string]int

// Member 是一个成员当前的状态。
type Member struct {
	Seq   int
	State map[string]string
}

type entry struct {
	member string
	seq    int
	state  map[string]string
}

// Sync 是队伍状态同步台账。
type Sync struct {
	members  map[string]Member
	buffered map[string]map[int]map[string]string
	log      []entry
	scanned  int
}

// New 建一份空的同步台账。
func New() *Sync {
	return &Sync{
		members:  map[string]Member{},
		buffered: map[string]map[int]map[string]string{},
	}
}

func clone(state map[string]string) map[string]string {
	out := map[string]string{}
	for key, value := range state {
		out[key] = value
	}
	return out
}

// Apply 应用某个成员的一次状态更新，返回这次更新是否当场生效。
func (s *Sync) Apply(member string, seq int, state map[string]string) bool {
	for range s.log {
		s.scanned++
	}
	s.log = append(s.log, entry{member: member, seq: seq, state: state})
	s.members[member] = Member{Seq: seq, State: clone(state)}
	return true
}

// Version 返回当前的版本向量。
func (s *Sync) Version() Vector {
	out := Vector{}
	for name, item := range s.members {
		out[name] = item.Seq
	}
	return out
}

// State 返回某个成员当前的状态。
func (s *Sync) State(member string) map[string]string {
	return clone(s.members[member].State)
}

// Pending 返回某个成员暂存待补的更新条数。
func (s *Sync) Pending(member string) int { return len(s.buffered[member]) }

// Diff 返回本地比对端新的成员，按名字升序。
func (s *Sync) Diff(peer Vector) []string {
	mine, theirs := 0, 0
	for _, seq := range s.Version() {
		if seq > mine {
			mine = seq
		}
	}
	for _, seq := range peer {
		if seq > theirs {
			theirs = seq
		}
	}
	if mine <= theirs {
		return []string{}
	}
	out := []string{}
	for name := range s.members {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}

// Missing 返回本地比对端落后的成员，按名字升序。
func (s *Sync) Missing(peer Vector) []string {
	mine, theirs := 0, 0
	for _, seq := range s.Version() {
		if seq > mine {
			mine = seq
		}
	}
	for _, seq := range peer {
		if seq > theirs {
			theirs = seq
		}
	}
	if theirs <= mine {
		return []string{}
	}
	out := []string{}
	for name := range peer {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}

// Scanned 是累计扫描的历史更新数（规模观测）。
func (s *Sync) Scanned() int { return s.scanned }
