package webui

import (
	"strings"
	"testing"
)

func TestDashboardEscapesAPIFieldsBeforeInnerHTML(t *testing.T) {
	checks := []string{
		"function html(value)",
		"html(safe(st.singbox_nodes))",
		"html(safe(st.subscription_count))",
		"html(safe(st.disabled_count))",
		"html(safe(st.subscription_total))",
		"html(p.protocol)",
		"html(p.exit_ip)",
		"html(sourceLabel(p))",
		"html(regionOf(p))",
		"html(sub.name)",
		"html(activeCount)",
		"html(pausedCount)",
		"html(disabledCount)",
		"html(line)",
		"html(s.session_id)",
		"html(region)",
		"function addressArg(address){return encodeURIComponent(String(address||'')).replace(/[!'()*]/g,c=>'%'+c.charCodeAt(0).toString(16).toUpperCase())}",
		"function proxyIDArg(proxy){const id=Number(proxy&&proxy.id);return Number.isFinite(id)?String(id):'0'}",
	}
	for _, check := range checks {
		t.Run(check, func(t *testing.T) {
			if !strings.Contains(dashboardHTML, check) {
				t.Fatalf("dashboardHTML missing escaped API-field usage %q", check)
			}
		})
	}

	for _, unsafe := range []string{
		"+safe(st.",
		"+safe(sub.",
		"+String(line).replace",
	} {
		t.Run("reject "+unsafe, func(t *testing.T) {
			if strings.Contains(dashboardHTML, unsafe) {
				t.Fatalf("dashboardHTML still contains unsafe innerHTML pattern %q", unsafe)
			}
		})
	}
}

func TestDashboardNodeStateAndRegionDistributionUseAvailableKnownRegions(t *testing.T) {
	checks := []string{
		// BUG-50: isAvailable 必须排除 user_paused，被暂停节点不算可用。
		"function isUserPaused(p){return !!(p&&(p.user_paused===true||Number(p.user_paused)===1))}",
		"function isAvailable(proxy){return !isUserPaused(proxy)&&(proxy.status==='active'||proxy.status==='degraded')&&Number(proxy.fail_count||0)<3}",
		// BUG-50: nodeState 主判据改为 user_paused（存储层新口径 status 仍为 active）。
		"function nodeState(p){if(isUserPaused(p)||p.status==='paused')return 'paused';if(isAvailable(p))return 'ok';if(p.status==='disabled'||Number(p.fail_count||0)>=3)return 'failed';return 'pending'}",
		// BUG-51: allRegions/地域分布均经 isAvailable 过滤，因此不再计入 user_paused 节点。
		"allRegions=Array.from(new Set(allProxies.filter(p=>isAvailable(p)&&isKnownRegion(p)).map(regionOf))).sort()",
		"allProxies.filter(p=>isAvailable(p)&&isKnownRegion(p)).forEach",
		"available nodes",
		"暂无可用地域数据",
	}
	for _, check := range checks {
		t.Run(check, func(t *testing.T) {
			if !strings.Contains(dashboardHTML, check) {
				t.Fatalf("dashboardHTML missing available-region invariant %q", check)
			}
		})
	}

	// 回归防护：isAvailable 必须真正读取 user_paused，不得退回旧口径（只看 status/fail_count）。
	if strings.Contains(dashboardHTML, "function isAvailable(proxy){return (proxy.status==='active'||proxy.status==='degraded')&&Number(proxy.fail_count||0)<3}") {
		t.Fatal("dashboardHTML isAvailable reverted to legacy form that ignores user_paused (BUG-50)")
	}
	// 回归防护：nodeState 不得只靠 status==='paused' 判定暂停（新数据 status 恒为 active，永不命中）。
	if strings.Contains(dashboardHTML, "function nodeState(p){if(p.status==='paused')return 'paused';") {
		t.Fatal("dashboardHTML nodeState reverted to legacy status==='paused' only check (BUG-50)")
	}
}

// TestDashboardPausedNodeTogglesToEnableButton 验证 BUG-50 修复后操作列按钮逻辑：
// user_paused 节点经 nodeState 归为 'paused'，renderProxies 据此显示“启用”按钮（enable=true），
// 非暂停节点显示“停用”按钮（enable=false）。
func TestDashboardPausedNodeTogglesToEnableButton(t *testing.T) {
	checks := []string{
		// toggleBtn 根据 nodeState 结果分支：paused -> 启用(true)，否则 停用(false)。
		"const toggleBtn=(st==='paused')?",
		",true)\">启用</button>'",
		",false)\">停用</button>'",
	}
	for _, check := range checks {
		t.Run(check, func(t *testing.T) {
			if !strings.Contains(dashboardHTML, check) {
				t.Fatalf("dashboardHTML missing paused-toggle-button invariant %q", check)
			}
		})
	}
}

func TestDashboardShowsPausedSubscriptionCountsAndLabels(t *testing.T) {
	checks := []string{
		"const pausedCount=Number(sub.paused_count??Math.max(0,proxyCount-activeCount-disabledCount))",
		"' 暂停 / '",
		"html(disabledCount)+' 不可用",
		"const badge=paused?'<span class=\"badge warn\">已暂停</span>'",
	}
	for _, check := range checks {
		t.Run(check, func(t *testing.T) {
			if !strings.Contains(dashboardHTML, check) {
				t.Fatalf("dashboardHTML missing paused subscription display %q", check)
			}
		})
	}
}

func TestDashboardShowsExplainedSingBoxStatusInsteadOfBoolFailure(t *testing.T) {
	checks := []string{
		"const status=String(st.singbox_status||(st.singbox_running?'running':'stopped'))",
		"no_tunnel_nodes:'无需运行'",
		"partial:'部分就绪'",
		"failed:'启动失败'",
		"<span class=\"k\">状态原因</span>",
		"html(reason)",
		"html(safe(st.singbox_ready_ports))+'/'+html(safe(st.singbox_total_ports))",
	}
	for _, check := range checks {
		t.Run(check, func(t *testing.T) {
			if !strings.Contains(dashboardHTML, check) {
				t.Fatalf("dashboardHTML missing explained sing-box status %q", check)
			}
		})
	}

	if strings.Contains(dashboardHTML, "running?'运行中':'未运行'") {
		t.Fatal("dashboardHTML still maps singbox_running=false directly to 未运行")
	}
}

func TestDashboardAsyncEntrypointsUseUnifiedErrorHandling(t *testing.T) {
	checks := []string{
		"async function runAsync(label, fn)",
		"try{data=JSON.parse(text)}catch(err){if(!res.ok)throw new Error(res.statusText||('HTTP '+res.status));throw new Error('响应解析失败')}",
		"async function refreshAll(){return runAsync('刷新失败'",
		"async function addManualNode(){return runAsync('添加失败'",
		"async function editManualRegion(id,address){return runAsync('地域更新失败'",
		"async function editManualNote(id,address){return runAsync('备注更新失败'",
		"async function deleteManualNode(id,address){return runAsync('删除失败'",
		"async function toggleProxy(id,address,enable){return runAsync('操作失败'",
		"async function addSubscription(){return runAsync('添加失败'",
		"async function refreshSub(id){return runAsync('刷新失败'",
		"async function refreshAllSubs(){return runAsync('刷新失败'",
		"async function toggleSub(id){return runAsync('切换失败'",
		"async function deleteSub(id){return runAsync('删除失败'",
		"async function saveConfig(){return runAsync('保存失败'",
		"setInterval(()=>runAsync('自动刷新失败'",
		"setInterval(()=>runAsync('日志刷新失败'",
	}
	for _, check := range checks {
		t.Run(check, func(t *testing.T) {
			if !strings.Contains(dashboardHTML, check) {
				t.Fatalf("dashboardHTML missing unified async handling %q", check)
			}
		})
	}

	if strings.Contains(dashboardHTML, ".catch(err=>showToast(err.message))") {
		t.Fatal("dashboardHTML still uses one-off catch instead of runAsync")
	}
}

func TestDashboardProxyActionsUseProxyIDAsPrimaryIdentity(t *testing.T) {
	checks := []string{
		"const id=proxyIDArg(p)",
		"toggleProxy('+id+',decodeURIComponent",
		"editManualRegion('+id+',decodeURIComponent",
		"editManualNote('+id+',decodeURIComponent",
		"deleteManualNode('+id+',decodeURIComponent",
		"const current=allProxies.find(p=>Number(p.id)===Number(id))||{}",
		"JSON.stringify({id,address,region})",
		"JSON.stringify({id,address,note})",
		"JSON.stringify({id,address})",
		"JSON.stringify({id,address,enable})",
	}
	for _, check := range checks {
		t.Run(check, func(t *testing.T) {
			if !strings.Contains(dashboardHTML, check) {
				t.Fatalf("dashboardHTML missing proxy-id action invariant %q", check)
			}
		})
	}

	for _, unsafe := range []string{
		"toggleProxy(decodeURIComponent",
		"editManualRegion(decodeURIComponent",
		"editManualNote(decodeURIComponent",
		"deleteManualNode(decodeURIComponent",
		"allProxies.find(p=>p.address===address)",
		"JSON.stringify({address,region})",
		"JSON.stringify({address,note})",
		"JSON.stringify({address})",
		"JSON.stringify({address,enable})",
	} {
		t.Run("reject "+unsafe, func(t *testing.T) {
			if strings.Contains(dashboardHTML, unsafe) {
				t.Fatalf("dashboardHTML still depends on address-only proxy action pattern %q", unsafe)
			}
		})
	}
}

func TestDashboardDoesNotShowOKBadgeForEmptySessionRegion(t *testing.T) {
	checks := []string{
		"const region=String(s.region||'').trim().toLowerCase()",
		"const regionBadge=region&&region!=='unknown'?'<span class=\"badge ok\">'+html(region).toUpperCase()+'</span> ':'<span class=\"badge gray\">未知</span> '",
		"'+regionBadge+masked+'",
	}
	for _, check := range checks {
		t.Run(check, func(t *testing.T) {
			if !strings.Contains(dashboardHTML, check) {
				t.Fatalf("dashboardHTML missing empty-session-region guard %q", check)
			}
		})
	}

	if strings.Contains(dashboardHTML, "'<span class=\"badge ok\">'+html(s.region)") {
		t.Fatal("dashboardHTML still renders an ok badge directly from s.region")
	}
}

func TestDashboardConnectionExampleAvoidsHttpbinSinglePoint(t *testing.T) {
	if strings.Contains(dashboardHTML, "httpbin.org") {
		t.Fatal("dashboardHTML still uses httpbin.org as a connection example target")
	}
	if !strings.Contains(dashboardHTML, "https://www.gstatic.com/generate_204") {
		t.Fatal("dashboardHTML missing stable HTTPS connection example target")
	}
}
