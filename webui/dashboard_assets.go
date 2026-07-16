package webui

// dashboard_assets.go 将 dashboard 的 CSS/JS 从 HTML 中分离为 Go 常量，
// 由 /assets/dashboard.css 与 /assets/dashboard.js 路由下发（带内容 hash 的 ETag、支持 304）。
// 仍为 Go 内嵌字符串，不落地独立文件、不引入前端构建链。

const dashboardCSS = `/* 后台样式；令牌对齐 docs/orbit-dashboard.html space/day */
.num{font-variant-numeric:tabular-nums;font-feature-settings:"tnum" 1}
.starfield{position:fixed;inset:0;z-index:0;pointer-events:none;overflow:hidden}
[data-theme="day"] .starfield{opacity:0}
.stars{position:absolute;inset:-50%;
 background-image:
  radial-gradient(1.4px 1.4px at 20% 30%,rgba(220,235,255,.9),transparent),
  radial-gradient(1.2px 1.2px at 70% 60%,rgba(180,215,255,.8),transparent),
  radial-gradient(1px 1px at 40% 80%,rgba(255,255,255,.75),transparent),
  radial-gradient(1.6px 1.6px at 85% 20%,rgba(190,220,255,.85),transparent);
 background-size:520px 520px,360px 360px,300px 300px,600px 600px;
 animation:star-drift 120s linear infinite}
.nebula{position:absolute;inset:0;
 background:
  radial-gradient(40% 30% at 22% 18%,rgba(59,141,255,.16),transparent 70%),
  radial-gradient(36% 28% at 82% 26%,rgba(124,196,255,.11),transparent 70%);
 filter:blur(6px)}
@keyframes star-drift{to{transform:translate(-260px,-260px)}}
@media (prefers-reduced-motion:reduce){.stars,.nebula{animation:none!important}}
/* 深空 space（默认） */
:root,[data-theme="space"]{
 --space-0:#04060e; --space-1:#080d1a; --space-2:#0e1626; --space-3:#152036;
 --bg:var(--space-0); --panel:rgba(14,22,38,.55); --panel-solid:#0e1626; --panel-2:rgba(20,32,54,.62);
 --soft:var(--space-3); --panel-3:var(--space-3); --hover:var(--space-3);
 --ink:#eaf1ff; --ink-2:#aebfd8; --muted:#8fa0bf; --ink-soft:var(--muted); --ink-3:var(--muted);
 --line:rgba(90,160,255,.16); --hairline:rgba(90,160,255,.09);
 --accent:#3b8dff; --accent-ink:#fff; --accent-strong:#7cc4ff; --accent-soft:rgba(59,141,255,.16); --signal:#3b8dff;
 --q-s:#7cc4ff; --q-a:#3b8dff; --q-b:#1f56c8; --q-c:#5a6480;
 --ok:#2fbf87; --warn:#f5b544; --danger:#ff5c7a; --gray:#6f7a96;
 --sun-core:#fff; --sun-halo:#9ccaff; --sun-energy:#3b8dff;
 --glow-s:0 0 12px rgba(124,196,255,.6),0 0 34px rgba(124,196,255,.28);
 --glow-a:0 0 12px rgba(59,141,255,.55),0 0 34px rgba(59,141,255,.24);
 --glow-b:0 0 12px rgba(31,86,200,.5),0 0 30px rgba(31,86,200,.2);
 --glow-c:0 0 10px rgba(90,100,128,.4),0 0 22px rgba(90,100,128,.16);
 --glow-ok:0 0 10px rgba(47,191,135,.6);
 --sh-sm:0 2px 8px rgba(0,0,0,.4);
 --sh-md:0 12px 40px rgba(0,0,0,.55),0 2px 8px rgba(0,0,0,.4);
 --sh-lg:0 30px 80px rgba(0,0,0,.65); --shadow:var(--sh-md);
 --radius:16px; --ease:cubic-bezier(.16,1,.3,1); --ease-out:var(--ease);
 --t-micro:150ms; --t-panel:280ms; --t-fast:150ms; --fs-caption:11px;
 --sidebar-w:236px; --sidebar-w-collapsed:74px; --topbar-h:60px;
 --accent-grad:linear-gradient(135deg,var(--q-s),var(--q-a) 55%,var(--q-b));
 --bg-canvas:radial-gradient(120% 120% at 50% 0%,#0b1226 0%,#070c18 46%,#04060e 100%);
 --surface-2:var(--space-3);
}
/* 白昼 day */
[data-theme="day"]{
 --space-0:#eef3fb; --space-1:#f6f9fe; --space-2:#fff; --space-3:#e9f0fb;
 --bg:var(--space-0); --panel:#fff; --panel-solid:#fff; --panel-2:#f4f8ff;
 --soft:#f4f8ff; --panel-3:var(--space-3); --hover:var(--space-3);
 --ink:#0f1b2e; --ink-2:#3a4a63; --muted:#5f6f8c; --ink-soft:var(--muted); --ink-3:var(--muted);
 --line:#e2e9f4; --hairline:rgba(30,60,110,.06);
 --accent:#1d6fe0; --accent-ink:#fff; --accent-strong:#1546a8; --accent-soft:rgba(29,111,224,.12); --signal:#3b8dff;
 --q-s:#4da3ff; --q-a:#1d6fe0; --q-b:#1546a8; --q-c:#8a93a6;
 --ok:#12a150; --warn:#c98a12; --danger:#e0485f; --gray:#9aa4ba;
 --sun-core:#fff; --sun-halo:#8fc0ff; --sun-energy:#1d6fe0;
 --glow-s:0 4px 14px rgba(77,163,255,.28); --glow-a:0 4px 14px rgba(29,111,224,.24);
 --glow-b:0 4px 12px rgba(21,70,168,.2); --glow-c:0 2px 8px rgba(138,147,166,.28);
 --glow-ok:0 3px 10px rgba(18,161,80,.3);
 --sh-sm:0 2px 6px rgba(30,60,110,.06);
 --sh-md:0 8px 30px rgba(30,60,110,.10),0 2px 6px rgba(30,60,110,.06);
 --sh-lg:0 24px 60px rgba(30,60,110,.16); --shadow:var(--sh-md);
 --radius:16px; --ease:cubic-bezier(.16,1,.3,1); --ease-out:var(--ease);
 --t-micro:150ms; --t-panel:280ms; --t-fast:150ms; --fs-caption:11px;
 --sidebar-w:236px; --sidebar-w-collapsed:74px; --topbar-h:60px;
 --accent-grad:linear-gradient(135deg,var(--q-s),var(--q-a) 55%,var(--q-b));
 --bg-canvas:radial-gradient(120% 120% at 50% 0%,#f2f7ff 0%,#e8f0fb 100%);
 --surface-2:var(--panel-2);
}
*{box-sizing:border-box}
body{margin:0;background:var(--bg-canvas,var(--bg));color:var(--ink);position:relative;z-index:1;
 font-family:"Segoe UI","PingFang SC","Microsoft YaHei",system-ui,sans-serif;font-size:14px;line-height:1.55;
 -webkit-font-smoothing:antialiased}
.sidebar,.main,.scrim{position:relative;z-index:1}
.sidebar{z-index:40}
.scrim{z-index:35}
button,input,select,textarea{font:inherit;color:inherit}
a{color:var(--accent)}
/* 统一可见焦点：绝不移除 outline，聚焦元素 2px 描边 */
:focus-visible{outline:2px solid var(--accent);outline-offset:2px;border-radius:4px}

/* ===== 布局骨架：左侧边栏 + 右主区 ===== */
.app{min-height:100vh}
.sidebar{position:fixed;top:0;left:0;bottom:0;width:var(--sidebar-w);z-index:40;
 background:linear-gradient(180deg,var(--panel-2),var(--panel));border-right:1px solid var(--line);
 display:flex;flex-direction:column;backdrop-filter:blur(14px);
 transition:width var(--t-micro) var(--ease),transform var(--t-panel) var(--ease)}
.sidebar.preload{transition:none}
.sidebar-brand{display:flex;align-items:center;gap:12px;height:var(--topbar-h);padding:0 18px;
 border-bottom:1px solid var(--hairline);overflow:hidden}
.mark{position:relative;flex:0 0 auto;width:38px;height:38px;border-radius:11px;
 background:radial-gradient(circle at 35% 30%,#dbeaff,var(--accent) 60%,var(--q-b,#1546a8));
 color:var(--accent-ink);display:grid;place-items:center;font-weight:900;font-size:13px;
 box-shadow:var(--glow-a,0 0 12px rgba(59,141,255,.4));z-index:1}
body.sidebar-collapsed .sidebar-foot .status-pill .lbl{opacity:0;width:0;overflow:hidden}
.sidebar-brand .bt{min-width:0;white-space:nowrap;overflow:hidden;text-overflow:ellipsis;font-weight:800;font-size:15px;
 transition:opacity var(--t-micro) var(--ease)}
.sidebar-nav{flex:1;padding:8px;display:flex;flex-direction:column;gap:4px;overflow-y:auto}
.navitem{display:flex;align-items:center;gap:12px;width:100%;padding:10px 12px;border-radius:10px;
 background:none;border:0;cursor:pointer;color:var(--muted);font-weight:600;text-align:left;white-space:nowrap;
 transition:background var(--t-micro) var(--ease),color var(--t-micro) var(--ease)}
.navitem .ico{flex:0 0 auto;width:20px;height:20px;display:grid;place-items:center}
.navitem .ico svg{width:20px;height:20px;display:block}
.navitem .lbl{min-width:0;overflow:hidden;text-overflow:ellipsis;transition:opacity var(--t-micro) var(--ease)}
.navitem:hover{background:color-mix(in srgb,var(--accent) 12%,transparent);color:var(--ink)}
.navitem.active{background:color-mix(in srgb,var(--accent) 16%,transparent);color:var(--ink);
 border:1px solid color-mix(in srgb,var(--accent) 40%,transparent);box-shadow:var(--glow-a,none)}
.navitem.active .ico svg{color:var(--accent);filter:drop-shadow(0 0 6px color-mix(in srgb,var(--accent) 70%,transparent))}
.sidebar-nav .lab{font-size:10px;letter-spacing:.16em;text-transform:uppercase;color:var(--muted);margin:12px 10px 4px;font-weight:700;white-space:nowrap}
body.sidebar-collapsed .sidebar-nav .lab{opacity:0;height:8px;margin:6px 0 0}
.sidebar-brand .bt small{display:block;font-size:10px;font-weight:600;color:var(--muted);letter-spacing:.14em;text-transform:uppercase}
.sidebar-foot{padding:12px;border-top:1px solid var(--hairline,var(--line));display:flex;flex-direction:column;gap:8px}
.sidebar-collapse{display:flex;align-items:center;gap:12px;width:100%;padding:9px 12px;border-radius:11px;
 border:1px solid var(--line);background:var(--panel-2);cursor:pointer;color:var(--ink-2);font-weight:600;text-align:left;white-space:nowrap}
.sidebar-collapse:hover{border-color:var(--accent);color:var(--accent)}
.sidebar-collapse .ico{flex:0 0 auto;width:19px;height:19px;display:grid;place-items:center;transition:transform var(--t-panel) var(--ease)}
.sidebar-collapse .ico svg{width:19px;height:19px;display:block}
.sidebar-collapse .lbl{min-width:0;overflow:hidden;text-overflow:ellipsis;transition:opacity var(--t-micro) var(--ease)}

/* 折叠态：仅图标，文字淡出隐藏（文字语义走 title + aria-label） */
body.sidebar-collapsed .sidebar{width:var(--sidebar-w-collapsed)}
body.sidebar-collapsed .navitem .lbl,
body.sidebar-collapsed .sidebar-collapse .lbl,
body.sidebar-collapsed .sidebar-brand .bt,
body.sidebar-collapsed .sidebar-foot .status-pill .lbl{opacity:0;width:0;pointer-events:none}
body.sidebar-collapsed .navitem,
body.sidebar-collapsed .sidebar-collapse{justify-content:center;gap:0;padding-left:0;padding-right:0}
body.sidebar-collapsed .sidebar-collapse .ico{transform:rotate(180deg)}

/* 主区随侧边栏宽度让位 */
.main{margin-left:var(--sidebar-w);transition:margin-left var(--t-micro) var(--ease)}
body.sidebar-collapsed .main{margin-left:var(--sidebar-w-collapsed)}

/* ===== 顶栏 ===== */
.topbar{position:sticky;top:0;z-index:30;height:var(--topbar-h);
 background:color-mix(in srgb,var(--panel-solid,var(--panel)) 82%,transparent);
 border-bottom:1px solid var(--line);display:flex;align-items:center;gap:12px;padding:0 20px;
 backdrop-filter:blur(14px)}
.topbar h1{margin:0;font-size:15px;font-weight:800;letter-spacing:.01em}
.hamburger{display:none;width:36px;height:36px;border:1px solid var(--line);background:var(--panel);
 color:var(--muted);border-radius:8px;cursor:pointer;align-items:center;justify-content:center;
 transition:border-color var(--t-micro) var(--ease),color var(--t-micro) var(--ease),background var(--t-micro) var(--ease),transform var(--t-micro) var(--ease)}
.hamburger:hover{border-color:var(--accent);color:var(--accent);background:color-mix(in srgb,var(--accent) 12%,var(--panel))}
.hamburger:active{transform:scale(.96)}
.topbar .status-pill{margin-left:4px}
.topbar-spacer{flex:1}
.overview-grid{display:grid;grid-template-columns:1.6fr 1fr;gap:16px;align-items:stretch;margin-bottom:16px}
.overview-side{display:flex;flex-direction:column;gap:16px}
.overview-side .card{margin-bottom:0}
.overview-side .conn{grid-template-columns:1fr}
.overview-side .cmd{font-size:12px;padding:12px}
.api-grid{display:grid;grid-template-columns:1.2fr 1fr;gap:16px;align-items:start}
@media(max-width:1100px){.overview-grid,.api-grid{grid-template-columns:1fr}}
.filters{display:flex;flex-wrap:wrap;gap:8px;align-items:center;margin-bottom:10px}
.input.sm{min-width:110px}
.input.mid{min-width:140px}
.input.narrow{width:110px}
.input.grow{flex:1;min-width:160px}
.page#page-settings .form-grid{display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:14px}
.page#page-settings .form-grid .field.full{grid-column:1/-1}
@media(max-width:800px){.page#page-settings .form-grid{grid-template-columns:1fr}}
.topbar .actions{display:flex;align-items:center;gap:8px}
.iconlink{display:inline-flex;align-items:center;justify-content:center;width:36px;height:36px;
 border:1px solid var(--line);border-radius:8px;color:var(--muted);text-decoration:none;background:var(--panel);
 transition:border-color var(--t-micro) var(--ease),color var(--t-micro) var(--ease),background var(--t-micro) var(--ease),transform var(--t-micro) var(--ease)}
.iconlink:hover{border-color:var(--accent);color:var(--accent);background:color-mix(in srgb,var(--accent) 12%,var(--panel))}
.iconlink:active{transform:scale(.96)}
.iconlink svg{width:20px;height:20px;display:block}

.btn{border:1px solid var(--line);background:var(--panel);border-radius:10px;padding:8px 14px;
 cursor:pointer;text-decoration:none;font-weight:600;white-space:nowrap;color:var(--ink);
 transition:border-color var(--t-micro) var(--ease),background var(--t-micro) var(--ease)}
.btn:hover{border-color:var(--accent)}
.btn.primary{background:var(--accent);border-color:var(--accent);color:var(--accent-ink)}
.btn.danger{color:var(--danger)}

.wrap{max-width:1320px;margin:0 auto;padding:24px}
.page{display:none}
.page.active{display:block}

/* ===== 遮罩（移动端抽屉） ===== */
.scrim{position:fixed;inset:0;background:rgba(12,18,30,.5);z-index:35;opacity:0;pointer-events:none;
 transition:opacity var(--t-panel) var(--ease)}
body.drawer-open .scrim{opacity:1;pointer-events:auto}

/* 指标卡 */
.metrics{display:grid;grid-template-columns:repeat(auto-fit,minmax(160px,1fr));gap:16px;margin-bottom:24px}
.metric{position:relative;overflow:hidden;background:linear-gradient(180deg,var(--panel-2),var(--panel));
 border:1px solid var(--line);border-radius:var(--radius);padding:16px 18px;box-shadow:var(--sh-md,var(--shadow));
 transition:transform var(--t-panel) var(--ease),border-color var(--t-panel)}
.metric:hover{transform:translateY(-2px);border-color:color-mix(in srgb,var(--accent) 34%,var(--line))}
.metric::before{content:"";position:absolute;left:0;top:0;height:3px;width:100%;
 background:linear-gradient(90deg,var(--q-s,var(--accent)),var(--q-a,var(--accent)),var(--q-b,var(--accent)),var(--q-c,var(--gray)))}
.metric .label{font-size:11px;letter-spacing:.06em;text-transform:uppercase;color:var(--muted);font-weight:700}
.metric .value{font-size:30px;font-weight:800;margin:6px 0 2px;letter-spacing:-.02em;
 text-shadow:0 0 20px color-mix(in srgb,var(--accent) 40%,transparent)}
.metric .note{font-size:11px;color:var(--muted)}

.card{background:linear-gradient(180deg,var(--panel-2),var(--panel));border:1px solid var(--line);
 border-radius:var(--radius);box-shadow:var(--sh-md,var(--shadow));margin-bottom:16px;overflow:hidden}
.card-head{display:flex;align-items:center;justify-content:space-between;gap:12px;padding:14px 18px;border-bottom:1px solid var(--hairline);flex-wrap:wrap}
.card-head h3{margin:0;font-size:14px;font-weight:800;letter-spacing:.01em}
.card-head .tools{display:flex;gap:8px;flex-wrap:wrap;align-items:center}
.card-body{padding:18px}
.two-col{display:grid;grid-template-columns:minmax(0,1fr) minmax(0,1fr);gap:16px;align-items:start}
@media(max-width:900px){.two-col{grid-template-columns:1fr}}

/* ===== 连接指引 ===== */
.conn{display:grid;grid-template-columns:repeat(auto-fit,minmax(224px,1fr));gap:16px}
.conn-item{background:var(--soft);border:1px solid var(--line);border-radius:12px;padding:16px}
.conn-item .k{font-size:11px;text-transform:uppercase;letter-spacing:.06em;color:var(--muted);font-weight:700}
.conn-item .v{font-family:"Consolas",monospace;font-size:15px;font-weight:700;margin-top:8px;word-break:break-all}
.conn-item .desc{font-size:12px;color:var(--muted);margin-top:4px}
.cmd{background:#0f1626;color:#cdd8ec;border-radius:10px;padding:16px;font-family:"Consolas",monospace;
 font-size:13px;overflow-x:auto;white-space:pre;margin-top:16px}
.cmd-hint{font-size:12px;color:var(--muted);line-height:1.7;margin-top:8px}
.cmd-hint code{font-family:"Consolas",monospace;font-size:12px;background:var(--soft);color:var(--accent);
 padding:1px 6px;border-radius:6px;border:1px solid var(--line)}
.cmd-hint b{color:var(--ink)}
.notice{display:flex;gap:8px;align-items:flex-start;background:var(--soft);border-left:3px solid var(--warn);
 border-radius:8px;padding:12px;font-size:13px;color:var(--muted);margin-top:16px}

.guide-row{display:flex;flex-wrap:wrap;gap:8px;font-family:"Consolas",monospace;font-size:13px;
 background:var(--soft);border-radius:8px;padding:8px 12px;margin-bottom:8px}
.guide-row b{color:var(--accent)}.guide-row span{color:var(--muted)}
.hint{font-size:12px;color:var(--muted);margin-top:8px}
.code-block{background:#0f1626;color:#cdd8ec;border-radius:10px;padding:16px;font-family:"Consolas",monospace;
 font-size:13px;overflow-x:auto;white-space:pre;margin:8px 0 0}

.toolbar{display:flex;flex-wrap:wrap;gap:8px;margin-bottom:16px}
.input{border:1px solid var(--line);border-radius:8px;padding:8px 12px;background:var(--panel);min-width:0;
 transition:border-color var(--t-micro) var(--ease)}
.input:focus{outline:none;border-color:var(--accent)}
.grow{flex:1;min-width:152px}

.table-wrap{width:100%;overflow-x:auto}
table{width:100%;border-collapse:collapse;font-size:13px}
th,td{padding:8px 10px;text-align:left;border-bottom:1px solid var(--line);white-space:nowrap;vertical-align:middle}
th{font-size:11px;text-transform:uppercase;letter-spacing:.04em;color:var(--muted);font-weight:700}
tbody tr:last-child td{border-bottom:none}
tbody tr:hover{background:color-mix(in srgb,var(--soft) 88%,transparent)}
.mono{font-family:"Consolas",monospace}
.empty{text-align:center;color:var(--muted);padding:24px 0}

.badge{display:inline-block;padding:2px 8px;border-radius:999px;font-size:11px;font-weight:700;background:var(--soft);color:var(--muted)}
.badge.ok{background:rgba(15,159,110,.14);color:var(--ok)}
.badge.blue{background:rgba(47,91,234,.13);color:var(--accent)}
.badge.warn{background:rgba(201,138,18,.15);color:var(--warn)}
.badge.danger{background:rgba(214,69,69,.14);color:var(--danger)}
.badge.gray{background:var(--soft);color:var(--gray)}
.muted{color:var(--muted)}

/* ===== 表头图标（CF / AI 统一图标语言：图标 + 短标签） ===== */
.th-ico{display:inline-flex;align-items:center;gap:5px;color:var(--muted);cursor:default}
.th-ico svg{width:15px;height:15px;display:block;flex:0 0 auto}
.th-ico .tx{font-size:var(--fs-caption,11px);font-weight:700;letter-spacing:.04em}

/* ===== AI 解锁标记（ChatGPT/Claude/Gemini/Grok：绿畅通 / 红阻断 / 灰未知） ===== */
.ai-marks{display:inline-flex;gap:3px;flex-wrap:wrap;vertical-align:middle;align-items:center}
.ai-mark{display:inline-grid;place-items:center;min-width:22px;height:18px;padding:0 5px;border-radius:5px;
 border:1px solid transparent;line-height:1;cursor:default;font-size:9px;font-weight:800;letter-spacing:.01em}
.ai-mark .nm{font-size:9px;font-weight:800;letter-spacing:.01em}
.ai-mark .gl{display:none}
.ai-mark.ok{color:var(--ok);border-color:color-mix(in srgb,var(--ok) 42%,transparent);background:color-mix(in srgb,var(--ok) 14%,transparent)}
.ai-mark.ok .nm{color:var(--ok)}
.ai-mark.bad{color:var(--danger);border-color:color-mix(in srgb,var(--danger) 42%,transparent);background:color-mix(in srgb,var(--danger) 12%,transparent)}
.ai-mark.bad .nm{color:var(--danger)}
.ai-mark.na{color:var(--gray);border-color:color-mix(in srgb,var(--gray) 40%,transparent);background:color-mix(in srgb,var(--gray) 10%,transparent);opacity:.75}
.ai-mark.na .nm{color:var(--gray)}

/* ===== 节点筛选栏：吸附 ===== */
.filter-toolbar{display:flex;flex-wrap:wrap;gap:8px;align-items:center;
 position:sticky;top:0;z-index:20;background:var(--panel);
 padding:12px 0;margin:0 0 12px;border-bottom:1px solid var(--line)}
.filter-toolbar .input{min-width:0}
.filter-toolbar .input.narrow{width:96px}
.filter-toolbar .input.mid{width:120px}
/* AI/CF 图标切换按钮（隐藏 select 承接原过滤语义） */
.hidden-select{position:absolute;width:1px;height:1px;padding:0;margin:-1px;overflow:hidden;clip:rect(0 0 0 0);border:0}
.filter-toggle{display:inline-flex;align-items:center;gap:6px;height:36px;padding:0 10px;border:1px solid var(--line);
 border-radius:8px;background:var(--panel);cursor:pointer;color:var(--muted);font-size:12px;font-weight:700;
 transition:border-color var(--t-micro) var(--ease),color var(--t-micro) var(--ease)}
.filter-toggle:hover{border-color:var(--accent)}
.filter-toggle[aria-pressed="true"]{border-color:var(--accent);color:var(--accent)}
.filter-toggle .ico{width:16px;height:16px;display:grid;place-items:center}
.filter-toggle .ico svg{width:16px;height:16px;display:block}
.filter-toggle .st{font-size:11px}
.check{display:inline-flex;align-items:center;gap:8px;font-size:12px;color:var(--muted);font-weight:600;cursor:pointer;user-select:none}
.check input{accent-color:var(--accent)}
.status-pill{display:inline-flex;align-items:center;gap:8px;font-size:12px;font-weight:700;color:var(--muted);letter-spacing:.04em}
.status-pill .dot{margin:0}

.mini{border:1px solid var(--line);background:var(--panel);border-radius:8px;padding:5px 10px;cursor:pointer;font-size:12px;font-weight:600;color:var(--ink);
 transition:border-color var(--t-micro) var(--ease)}
.mini:hover{border-color:var(--accent)}
.mini.primary{background:var(--accent);border-color:var(--accent);color:var(--accent-ink)}
.mini.danger{color:var(--danger)}

.region-row{display:flex;align-items:center;gap:12px;padding:8px 0}
.region-row strong{width:40px;font-size:13px}
.bar{flex:1;height:8px;background:var(--soft);border-radius:999px;overflow:hidden}
.bar span{display:block;height:100%;background:var(--accent);border-radius:999px}
.region-row .cnt{width:40px;text-align:right;color:var(--muted);font-size:13px}

.kv{display:flex;align-items:center;justify-content:space-between;gap:12px;padding:8px 0;border-bottom:1px solid var(--line);font-size:13px}
.kv:last-child{border-bottom:none}
.kv .k{color:var(--muted)}.kv .v{font-weight:700}
.dot{display:inline-block;width:8px;height:8px;border-radius:50%;margin-right:8px;vertical-align:middle}
.dot.on{background:var(--ok)}.dot.off{background:var(--danger)}.dot.warn{background:var(--warn)}.dot.idle{background:var(--gray)}

.sub-item{display:flex;align-items:center;justify-content:space-between;gap:12px;padding:14px;border-radius:12px;
 border:1px solid var(--hairline,var(--line));background:var(--panel-2);margin-bottom:10px;flex-wrap:wrap}
.sub-item .meta{min-width:0}
.sub-item .meta strong{font-size:14px;display:inline-flex;align-items:center;gap:8px}
.sub-item .meta .muted{margin-top:4px;font-size:12px;color:var(--muted)}
.mini-actions{display:flex;gap:6px;flex-wrap:wrap}

/* 会话页：可展开卡片（对齐设计稿） */
.session-list,.session-grid{display:flex;flex-direction:column;gap:12px}
.session-card{padding:0;border-radius:14px;border:1px solid var(--hairline,var(--line));background:var(--panel-2);overflow:hidden;
 transition:border-color var(--t-micro),box-shadow var(--t-micro)}
.session-card.open{border-color:color-mix(in srgb,var(--accent) 40%,var(--line));box-shadow:var(--glow-a,none)}
.session-card .head{display:flex;align-items:center;gap:12px;padding:14px 16px;cursor:pointer;user-select:none}
.session-card .head:hover{background:color-mix(in srgb,var(--accent) 6%,transparent)}
.session-card .sid{font-family:"Consolas",monospace;font-weight:800;font-size:14px;min-width:0;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.session-card .ttl{color:var(--ok);font-weight:800;font-size:12px;white-space:nowrap}
.session-card .ttl.warn{color:var(--warn)}.session-card .ttl.danger{color:var(--danger)}
.session-card .chips{display:flex;align-items:center;gap:6px;flex-wrap:wrap;flex:1;min-width:0}
.session-card .chev{flex:0 0 auto;width:22px;height:22px;display:grid;place-items:center;color:var(--muted);
 transition:transform var(--t-panel) var(--ease)}
.session-card.open .chev{transform:rotate(90deg);color:var(--accent)}
.session-card .body{display:none;padding:0 16px 16px;border-top:1px solid var(--hairline,var(--line))}
.session-card.open .body{display:block}
.session-card .detail-grid{display:grid;grid-template-columns:repeat(auto-fit,minmax(180px,1fr));gap:10px 16px;padding-top:14px}
.session-card .di{display:flex;flex-direction:column;gap:3px;min-width:0}
.session-card .di .k{font-size:10px;font-weight:700;letter-spacing:.06em;text-transform:uppercase;color:var(--muted)}
.session-card .di .v{font-size:13px;font-weight:600;color:var(--ink);word-break:break-all}
.session-card .di .v.mono{font-family:"Consolas",monospace;font-weight:700}
.ops{display:inline-flex;gap:4px;flex-wrap:nowrap;white-space:nowrap}
.toolbar.filters{padding:10px 12px;border:1px solid var(--hairline,var(--line));border-radius:12px;
 background:color-mix(in srgb,var(--panel-2) 80%,transparent);justify-content:flex-start;gap:8px;margin-bottom:10px}
/* 日志框：按视口拉高，吃掉底部空白 */
#page-logs .card{display:flex;flex-direction:column;min-height:calc(100vh - var(--topbar-h,60px) - 48px)}
#page-logs .card-body{flex:1;display:flex;flex-direction:column;min-height:0}
.logs{background:#060a14;border:1px solid var(--line);border-radius:12px;padding:12px 14px;overflow:auto;
 font-family:"Consolas","Cascadia Mono",monospace;font-size:12px;line-height:1.55;color:#c8d7ef;
 height:calc(100vh - var(--topbar-h,60px) - 180px);min-height:420px;max-height:none}
[data-theme="day"] .logs{background:#0f1726;color:#d7e4ff}
.log-line{padding:2px 4px;border-radius:4px;white-space:pre-wrap;word-break:break-all}
.log-line:hover{background:rgba(59,141,255,.08)}
.form-grid{display:grid;grid-template-columns:1fr 1fr;gap:14px}
@media(max-width:800px){.form-grid{grid-template-columns:1fr}}
.field{display:flex;flex-direction:column;gap:6px}
.field.full{grid-column:1/-1}
.field label{font-size:11px;font-weight:700;letter-spacing:.04em;text-transform:uppercase;color:var(--muted)}
.field .fh{font-size:11px;color:var(--muted)}
.field input,.field select,.field textarea{border:1px solid var(--line);border-radius:10px;padding:8px 12px;background:var(--panel-2);color:var(--ink)}
.legend{display:flex;gap:16px;flex-wrap:wrap;font-size:12px;color:var(--muted);margin-bottom:12px}
.legend span{display:flex;align-items:center;gap:6px}

/* ===== 骨架墓碑（shimmer） ===== */
.skeleton{display:block;background:linear-gradient(90deg,var(--soft) 25%,color-mix(in srgb,var(--line) 60%,var(--soft)) 37%,var(--soft) 63%);
 background-size:400% 100%;border-radius:8px;animation:sk-shimmer 1.4s ease infinite}
.sk-line{height:12px;margin:8px 0}
.sk-row{height:40px;margin:8px 0;border-radius:10px}
@keyframes sk-shimmer{0%{background-position:100% 0}100%{background-position:0 0}}
.skeleton-wrap{padding:8px 0}

.modal{position:fixed;inset:0;background:rgba(12,18,30,.5);display:none;align-items:flex-start;justify-content:center;padding:44px 16px;z-index:60;overflow:auto}
.modal.show{display:flex}
.dialog{background:var(--panel);border-radius:var(--radius);width:min(560px,100%);padding:24px;box-shadow:0 30px 80px rgba(10,16,30,.4)}
.dialog h3{margin:0 0 16px}
.form-grid{display:grid;grid-template-columns:1fr 1fr;gap:16px}
.field{display:flex;flex-direction:column;gap:8px}
.field.full{grid-column:1 / -1}
.field label{font-size:12px;color:var(--muted);font-weight:600}
.field input,.field select,.field textarea{border:1px solid var(--line);border-radius:8px;padding:8px 12px;background:var(--panel);width:100%}
.field textarea{min-height:120px;resize:vertical;font-family:"Consolas",monospace}
.field .fh{font-size:11px;color:var(--muted)}
.dialog-actions{display:flex;justify-content:flex-end;gap:8px;margin-top:24px}
.apikey-section{margin-top:24px;border-top:1px solid var(--line);padding-top:16px}
.apikey-section h4{margin:0 0 8px}

.toast{position:fixed;left:50%;bottom:24px;transform:translateX(-50%) translateY(20px);background:var(--ink);
 color:var(--bg);padding:12px 20px;border-radius:999px;font-weight:600;opacity:0;pointer-events:none;
 transition:opacity var(--t-panel) var(--ease),transform var(--t-panel) var(--ease);z-index:70}
.toast.show{opacity:1;transform:translateX(-50%) translateY(0)}

/* 节点分布图（rAF 动画） */
.orbit-card .card-body{display:flex;flex-direction:column;gap:12px}
.orbit-stage{position:relative;width:100%;aspect-ratio:16/9;margin:0 auto;overflow:hidden;border-radius:12px;
 background:radial-gradient(ellipse at 50% 45%,color-mix(in srgb,var(--accent) 10%,transparent),transparent 62%)}
.orbit-svg{position:absolute;inset:0;width:100%;height:100%;overflow:visible;pointer-events:none}
.orbit-ring{fill:none;stroke-width:1.2}
.orbit-beam{stroke:none;fill-rule:nonzero;pointer-events:none}
.orbit-beam-glow{stroke:none;fill-rule:nonzero;pointer-events:none;mix-blend-mode:screen}
.orbit-sats{position:absolute;inset:0;pointer-events:none}
.orbit-sat{position:absolute;transform:translate(-50%,-50%);pointer-events:auto;cursor:default;will-change:left,top,width,height}
.orbit-sat .ball{width:100%;height:100%;border-radius:50%;display:grid;place-items:center;
 background:radial-gradient(circle at 34% 28%,#fff 0%,var(--qc,var(--q-a)) 42%,color-mix(in srgb,var(--qc,var(--q-a)) 55%,#04060e) 100%);
 box-shadow:inset 0 -3px 8px rgba(0,0,0,.35),inset 2px 2px 6px rgba(255,255,255,.35),0 0 10px color-mix(in srgb,var(--qc,var(--q-a)) 45%,transparent);
 border:1px solid color-mix(in srgb,var(--qc,var(--q-a)) 60%,transparent)}
.orbit-sat .cc{font-size:10px;font-weight:900;color:#fff;text-shadow:0 1px 3px rgba(0,0,0,.6);letter-spacing:.02em}
.orbit-sat .cnt{position:absolute;top:-5px;right:-5px;min-width:13px;height:13px;padding:0 3px;border-radius:7px;
 display:grid;place-items:center;font-size:8px;font-weight:800;color:var(--ink);
 background:var(--panel-solid,var(--panel));border:1px solid var(--qc,var(--q-a))}
.orbit-sat.live .ball{box-shadow:inset 0 -3px 8px rgba(0,0,0,.35),inset 2px 2px 6px rgba(255,255,255,.4),0 0 16px color-mix(in srgb,var(--qc,var(--q-a)) 75%,transparent)}
.orbit-sun{position:absolute;left:50%;top:50%;width:84px;height:84px;transform:translate(-50%,-50%);border-radius:50%;
 display:grid;place-items:center;z-index:5;
 background:radial-gradient(circle at 42% 36%,#fff,var(--q-s,#7cc4ff) 46%,var(--q-b,#1f56c8) 78%,#0e2a5e);
 box-shadow:0 0 28px rgba(124,196,255,.65),0 0 60px rgba(59,141,255,.35),inset 0 0 18px rgba(255,255,255,.45)}
.orbit-sun-ring{position:absolute;inset:-12px;border-radius:50%;z-index:-1;
 background:conic-gradient(from 0deg,transparent,var(--accent),transparent 30%,var(--q-s,var(--accent)),transparent 60%,var(--accent),transparent);
 filter:blur(3px);animation:orbit-ring-spin 8s linear infinite}
.orbit-sun-halo{position:absolute;inset:-4px;border-radius:50%;z-index:-1;
 background:radial-gradient(circle,color-mix(in srgb,var(--accent) 35%,transparent),transparent 70%);animation:orbit-pulse 3s ease-in-out infinite}
.orbit-sun-lbl{text-align:center;color:#fff;text-shadow:0 1px 6px rgba(0,0,0,.5);z-index:1;line-height:1.2}
.orbit-sun-lbl .t{font-size:9px;letter-spacing:.14em;font-weight:700;opacity:.85;text-transform:uppercase}
.orbit-sun-lbl .ip{font-size:11px;font-weight:800;letter-spacing:.02em}
.orbit-legend{display:flex;flex-direction:column;align-items:center;gap:8px;font-size:11px;color:var(--muted)}
.orbit-legend-row{display:flex;flex-wrap:wrap;gap:14px;justify-content:center;align-items:center}
.orbit-legend b{display:inline-flex;align-items:center;gap:6px;font-weight:600;color:var(--ink-2,var(--ink))}
.orbit-legend .qd{width:10px;height:10px;border-radius:50%}
.orbit-legend .qd.s{background:var(--q-s);box-shadow:0 0 8px color-mix(in srgb,var(--q-s) 50%,transparent)}
.orbit-legend .qd.a{background:var(--q-a);box-shadow:0 0 8px color-mix(in srgb,var(--q-a) 50%,transparent)}
.orbit-legend .qd.b{background:var(--q-b);box-shadow:0 0 8px color-mix(in srgb,var(--q-b) 40%,transparent)}
.orbit-legend .qd.c{background:var(--q-c);box-shadow:0 0 6px color-mix(in srgb,var(--q-c) 35%,transparent)}
.beam-swatch{width:16px;height:3px;background:var(--q-a);display:inline-block;border-radius:2px;box-shadow:var(--glow-a)}
.wind-swatch{width:18px;height:8px;display:inline-block;background:linear-gradient(90deg,transparent,rgba(156,202,255,.7),transparent);border-radius:4px;filter:blur(.5px);box-shadow:0 0 8px rgba(124,196,255,.35)}
.lens-swatch{width:14px;height:14px;border-radius:50%;border:1.5px solid rgba(180,210,255,.55);box-shadow:inset 0 0 6px rgba(180,210,255,.35);display:inline-block}
.orbit-wind-plume{fill:url(#orbitWindPlume);pointer-events:none;mix-blend-mode:screen;filter:blur(1.2px)}
.orbit-wind-stream{fill:none;stroke-linecap:round;pointer-events:none;mix-blend-mode:screen;filter:blur(.6px)}
.orbit-wind-stream-core{fill:none;stroke-linecap:round;pointer-events:none;mix-blend-mode:screen}
.orbit-lens-halo{fill:url(#orbitLensFill);pointer-events:none;mix-blend-mode:screen}
.orbit-lens-rim{fill:none;stroke:rgba(180,210,255,.35);stroke-width:1.2;pointer-events:none}
@keyframes orbit-ring-spin{to{transform:rotate(360deg)}}
@keyframes orbit-pulse{0%,100%{opacity:.45}50%{opacity:.95}}

/* ===== 响应式：≤900px 侧边栏改 off-canvas 抽屉 ===== */
@media(max-width:900px){
 .sidebar{transform:translateX(-100%)}
 body.drawer-open .sidebar{transform:translateX(0);width:var(--sidebar-w)}
 body.drawer-open .navitem .lbl,
 body.drawer-open .sidebar-foot .btn .lbl,
 body.drawer-open .sidebar-brand .bt{opacity:1;width:auto}
 .main{margin-left:0}
 body.sidebar-collapsed .main{margin-left:0}
 .hamburger{display:inline-flex}
 .sidebar-toggle{display:none}
 .sidebar-collapse{display:none}
 .wrap{padding:16px}
}

/* ===== 尊重减少动效偏好：关闭非必要动画 ===== */
@media (prefers-reduced-motion:reduce){
 *{animation-duration:.001ms !important;animation-iteration-count:1 !important;transition-duration:.001ms !important}
 .orbit-sun-ring,.orbit-sun-halo,.skeleton{animation:none}
}
/* 筛选切换按钮内的服务名与状态点 */
.filter-toggle .nm{font-weight:700;color:var(--ink)}
.filter-toggle[aria-pressed="true"] .nm{color:var(--accent)}
.fdot{width:8px;height:8px;border-radius:50%;background:var(--gray);flex:0 0 auto}
.filter-toggle[aria-pressed="true"] .fdot{background:var(--accent)}`

const dashboardJS = `let allProxies=[];let allRegions=[];let configCache=null;let publicIP='';let orbitSessions=[];let gatewayCC='';
const PAGE_TITLES={overview:'总览',nodes:'节点',subs:'订阅',sessions:'会话',logs:'日志',api:'API',settings:'设置'};
function switchTab(name){document.querySelectorAll('.navitem').forEach(t=>t.classList.toggle('active',t.dataset.tab===name));document.querySelectorAll('.page').forEach(p=>p.classList.toggle('active',p.id==='page-'+name));const title=document.getElementById('pageTitle');if(title)title.textContent=PAGE_TITLES[name]||name;if(name==='settings'){runAsync('打开设置失败',async()=>{if(!configCache)await loadConfig();await loadAPIKeys()})}try{markViewLazy(name)}catch(e){}closeDrawer()}
function showToast(msg){const el=document.getElementById('toast');el.textContent=msg;el.classList.add('show');setTimeout(()=>el.classList.remove('show'),2600)}
async function api(path, options){const res=await fetch(path, Object.assign({headers:{'Content-Type':'application/json'}}, options||{}));if(res.status===401){location.href='/login';return null}const text=await res.text();let data={};if(text){try{data=JSON.parse(text)}catch(err){if(!res.ok)throw new Error(res.statusText||('HTTP '+res.status));throw new Error('响应解析失败')}}if(!res.ok)throw new Error(data.error||res.statusText||('HTTP '+res.status));return data}
function safe(value){return value===undefined||value===null||value===''?'--':String(value)}
function html(value){return safe(value).replace(/[&<>"']/g,c=>({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[c]))}
function errorMessage(err){return err&&err.message?err.message:String(err||'操作失败')}
async function runAsync(label, fn){try{return await fn()}catch(err){showToast((label?label+'：':'')+errorMessage(err));return null}}
async function logout(){return runAsync('退出失败',async()=>{const res=await fetch('/logout',{method:'POST'});if(!res.ok)throw new Error(res.statusText||('HTTP '+res.status));location.href='/login'})}
function refreshLater(){setTimeout(()=>runAsync('刷新失败',()=>Promise.all([loadSubscriptions(),loadStats(),loadProxies()])),4000)}
function maskAddress(address){if(!address)return '--';const parts=String(address).split(':');const host=parts[0]||address;if(host.length<=8)return host+(parts[1]?':'+parts[1]:'');return host.slice(0,4)+'...'+host.slice(-4)+(parts[1]?':'+parts[1]:'')}
function addressArg(address){return encodeURIComponent(String(address||'')).replace(/[!'()*]/g,c=>'%'+c.charCodeAt(0).toString(16).toUpperCase())}
function proxyIDArg(proxy){const id=Number(proxy&&proxy.id);return Number.isFinite(id)?String(id):'0'}
function regionOf(proxy){const region=String((proxy&&proxy.region)||'').trim().toLowerCase();return region||'unknown'}
function isKnownRegion(proxy){const region=regionOf(proxy);return region&&region!=='unknown'}
function isUserPaused(p){return !!(p&&(p.user_paused===true||Number(p.user_paused)===1))}
function isAvailable(proxy){return !isUserPaused(proxy)&&(proxy.status==='active'||proxy.status==='degraded')&&Number(proxy.fail_count||0)<3}
function stripColon(port){return String(port||'').replace(/^:/,'')}
async function refreshAll(){return runAsync('刷新失败',async()=>{await Promise.all([loadStats(),loadProxies(),loadSubscriptions(),loadConfig(),loadSessions(),loadLogs(),loadCustomStatus()]);loadPublicIP();showToast('数据已刷新')})}
async function loadCustomStatus(){const st=await api('/api/custom/status');if(!st)return;const box=document.getElementById('singbox-status');if(!box)return;const status=String(st.singbox_status||(st.singbox_running?'running':'stopped'));const reason=String(st.singbox_reason||status);const statusText={no_tunnel_nodes:'无需运行',running:'运行中',stopped:'已停止',partial:'部分就绪',failed:'启动失败'}[status]||status;const dotClass={no_tunnel_nodes:'idle',running:'on',stopped:'idle',partial:'warn',failed:'off'}[status]||'idle';const dot='<span class="dot '+dotClass+'"></span>';box.innerHTML='<div class="kv"><span>'+dot+'sing-box 引擎</span><span class="v">'+html(statusText)+'</span></div>'+'<div class="kv"><span class="k">状态原因</span><span class="v">'+html(reason)+'</span></div>'+'<div class="kv"><span class="k">转换节点</span><span class="v">'+html(safe(st.singbox_nodes))+'</span></div>'+'<div class="kv"><span class="k">端口就绪</span><span class="v">'+html(safe(st.singbox_ready_ports))+'/'+html(safe(st.singbox_total_ports))+'</span></div>'+'<div class="kv"><span class="k">订阅可用</span><span class="v">'+html(safe(st.subscription_count))+'</span></div>'+'<div class="kv"><span class="k">暂停/不可用节点</span><span class="v">'+html(safe(st.disabled_count))+'</span></div>'+'<div class="kv"><span class="k">订阅总数</span><span class="v">'+html(safe(st.subscription_total))+'</span></div>'}
function normalizeTheme(theme){const t=String(theme||'');if(t==='day'||t==='light')return 'day';return 'space'}
function applyTheme(theme){const t=normalizeTheme(theme);document.documentElement.setAttribute('data-theme',t);try{localStorage.setItem('gg-theme',t)}catch(e){}const btn=document.getElementById('theme-toggle');if(btn){btn.title=t==='space'?'切换到浅色':'切换到深色';btn.setAttribute('aria-label',btn.title)}try{if(document.getElementById('orbit-stage'))renderOrbitSystem()}catch(e){}}
function toggleTheme(){const cur=normalizeTheme(document.documentElement.getAttribute('data-theme'));applyTheme(cur==='space'?'day':'space')}
(function(){let t='space';try{t=localStorage.getItem('gg-theme')||'space'}catch(e){}applyTheme(t)})();
async function loadStats(){const stats=await api('/api/stats');if(!stats)return;document.getElementById('stat-total').textContent=safe(stats.total);document.getElementById('stat-http').textContent=safe(stats.http);document.getElementById('stat-socks5').textContent=safe(stats.socks5);document.getElementById('stat-subscription').textContent=safe(stats.subscription_count);document.getElementById('stat-sessions').textContent=safe(stats.active_sessions)}
async function loadProxies(){const data=await api('/api/proxies');if(!data)return;allProxies=Array.isArray(data)?data:[];allRegions=Array.from(new Set(allProxies.filter(p=>isAvailable(p)&&isKnownRegion(p)).map(regionOf))).sort();renderRegionFilter();renderProxies();renderRegions();renderOrbitSystem()}
function renderRegionFilter(){const select=document.getElementById('region-filter');const current=select.value;select.innerHTML='<option value="">全部地域</option>'+allRegions.map(r=>'<option value="'+html(r)+'">'+html(r).toUpperCase()+'</option>').join('');select.value=allRegions.includes(current)?current:''}
function sourceLabel(p){if(p.source==='manual')return '手工';return p.subscription_name?p.subscription_name:'订阅';}
function nodeLabel(p){if(p.source==='manual')return maskAddress(p.address);if(p.note)return p.note;return p.subscription_name?p.subscription_name:'订阅节点';}
// 节点状态：user_paused→已停用；可用→ok；disabled 且从未验证→pending；否则不可用。
// last_check 空/零值视为未验证（待验证），有 last_check 或 fail_count≥3 视为验证失败。
function hasLastCheck(p){const v=p&&p.last_check;if(v==null||v===''||v===false)return false;const s=String(v);if(s.indexOf('0001-')===0||s.indexOf('1970-01-01')===0)return false;return true}
function nodeState(p){if(isUserPaused(p)||p.status==='paused')return 'paused';if(isAvailable(p))return 'ok';if(Number(p.fail_count||0)>=3)return 'failed';if(p.status==='disabled')return hasLastCheck(p)?'failed':'pending';return 'pending'}
function stateBadge(st){switch(st){case 'ok':return '<span class="badge ok">可用</span>';case 'paused':return '<span class="badge warn">已停用</span>';case 'failed':return '<span class="badge danger">不可用</span>';default:return '<span class="badge gray">待验证</span>'}}
// abuserBadge: ipapi.is abuser_score<0 显示 "--"（未探测/查询失败）；否则显示 0.00-1.00 两位小数 + 颜色。
// 阈值：<0.10 绿(ok)、0.10-0.50 黄(warn)、>0.50 红(danger)。两源分开展示，不与 ip-api 聚合。
function abuserBadge(score){const n=Number(score);if(!Number.isFinite(n)||n<0)return '<span class="muted">--</span>';const cls=n<0.1?'ok':(n<=0.5?'warn':'danger');return '<span class="badge '+cls+'">'+html(n.toFixed(2))+'</span>'}
// ipapiFlagsBadges: ip-api 命中标记逗号串。proxy 红、hosting 黄、mobile 灰；seen=true 且无命中显"干净"绿；未探测显 "--"。
function ipapiFlagsBadges(flags,seen){const raw=String(flags||'').trim();if(raw===''){return seen?'<span class="badge ok">干净</span>':'<span class="muted">--</span>'}const cls={proxy:'danger',hosting:'warn',mobile:'gray'};return raw.split(',').map(f=>f.trim()).filter(Boolean).map(f=>'<span class="badge '+(cls[f]||'gray')+'">'+html(f)+'</span>').join(' ')}
// cfBadge: cf_blocked==1 显"拦截"红、==0 显"正常"绿、其它(-1/未探测)显 "--"。
function cfBadge(v){ v=Number(v); if(v===1)return '<span class="badge danger">拦截</span>'; if(v===0)return '<span class="badge ok">正常</span>'; return '<span class="muted">--</span>' }
// aiBadges: 解析 ai_reachability JSON（形如 {"openai":0,"claude":1,"grok":-1,"gemini":0}），
// 四服务正规短标签：ChatGPT/Claude/Grok/Gemini。0=畅通(绿)、1=阻断(红)、其它=未知(灰)。
// title 带全名与状态。空/非法 JSON 整体显 "--"。
function aiBadges(v){ const raw=String(v||'').trim(); if(raw===''){return '<span class="muted">--</span>'} let m; try{m=JSON.parse(raw)}catch(e){return '<span class="muted">--</span>'} if(!m||typeof m!=='object'){return '<span class="muted">--</span>'} const defs=[['openai','ChatGPT','ChatGPT'],['claude','Claude','Claude'],['grok','Grok','Grok'],['gemini','Gemini','Gemini']]; return '<span class="ai-marks">'+defs.map(function(d){const k=d[0],ab=d[1],full=d[2];const n=Number(m[k]);const cls=n===0?'ok':(n===1?'bad':'na');const glyph=n===0?'✓':(n===1?'✗':'–');const title=full+(n===0?' 畅通':(n===1?' 阻断':' 未知'));return '<span class="ai-mark '+cls+'" title="'+html(title)+'"><span class="nm">'+ab+'</span><span class="gl">'+glyph+'</span></span>'}).join('')+'</span>' }
function aiStateOf(p,svc){const raw=String((p&&p.ai_reachability)||'').trim();if(!raw)return 'unprobed';let m;try{m=JSON.parse(raw)}catch(e){return 'unprobed'}if(!m||typeof m!=='object')return 'unprobed';const n=Number(m[svc]);if(n===0)return 'unlocked';if(n===1)return 'blocked';return 'unprobed'}
function cfStateOf(p){const v=Number(p&&p.cf_blocked);if(v===0)return 'unlocked';if(v===1)return 'blocked';return 'unknown'}
function qualityOf(p){return String((p&&p.quality_grade)||'').trim().toUpperCase()}
function filterVal(id){const el=document.getElementById(id);return el?String(el.value||'').trim():''}
// starBtn: 星标切换按钮，★ 已加星 / ☆ 未加星。
function starBtn(p){ const id=proxyIDArg(p); const on=!!(p.starred===true||Number(p.starred)===1); return '<button class="mini" onclick="toggleStar('+id+','+(on?'true':'false')+')" title="星标">'+(on?'★':'☆')+'</button>' }
// randSession: 随机 6 位字母数字，用于复制凭据的 session 段。
function randSession(){ const cs='abcdefghijklmnopqrstuvwxyz0123456789'; let s=''; for(let i=0;i<6;i++)s+=cs[Math.floor(Math.random()*cs.length)]; return s }
// isDualProtocol: 节点是否为 sing-box mixed 入站(单端口同时服务 SOCKS5+HTTP)。
// 读存储层显式下发的 dual_protocol 字段,而非靠地址长相猜测——手动本机 direct socks5 节点
// 地址同为回环但只支持单协议,只有此显式标记能可靠区分。
function isDualProtocol(p){return !!(p&&(p.dual_protocol===true||Number(p.dual_protocol)===1))}
// protocolBadges: 协议列徽章。dual_protocol 节点(mixed 入站)渲染 SOCKS5+HTTP 两个徽章;
// 其余节点按存储的单一 protocol 渲染一个徽章(沿用 html 转义)。
function protocolBadges(p){ if(isDualProtocol(p))return '<span class="badge blue">SOCKS5</span> <span class="badge blue">HTTP</span>'; return '<span class="badge blue">'+html(p.protocol).toUpperCase()+'</span>' }
// isGatewayNode: dual_protocol(mixed 隧道)或回环本地地址必须经网关 DSL 连接；其余为可直连上游。
function isGatewayNode(p){if(isDualProtocol(p))return true;const a=String((p&&p.address)||'');return a.indexOf('127.0.0.1:')===0||a.indexOf('[::1]:')===0||a.indexOf('localhost:')===0}
function isDirectNode(p){return !isGatewayNode(p)}
// copyProxyCred: 直连节点复制 protocol://host:port（无网关密码）；网关节点复制 DSL 凭据到公网入口。
// 用户名/密码编码为 URL userinfo。成功 toast 不回显含真实密码的完整 URL。
function encodeProxyUserInfo(value){return encodeURIComponent(String(value||'')).replace(/[!'()*]/g,c=>'%'+c.charCodeAt(0).toString(16).toUpperCase())}
function copyProxyCred(id){ const p=allProxies.find(x=>Number(x.id)===Number(id)); if(!p)return; const addr=String(p.address||''); const scheme=isDualProtocol(p)?(confirm('确定复制 SOCKS5？取消则复制 HTTP')?'socks5':'http'):String(p.protocol||'socks5'); if(isDirectNode(p)){ const url=scheme+'://'+addr; navigator.clipboard.writeText(url).then(()=>showToast('已复制直连地址')).catch(()=>showToast('复制失败')); return } const base=(configCache&&configCache.proxy_auth_username)?configCache.proxy_auth_username:'username'; const region=regionOf(p); const user=base+'-region-'+region+'-session-'+randSession(); const rawPass=(configCache&&configCache.proxy_auth_password)?configCache.proxy_auth_password:''; const pass=rawPass||'PASSWORD'; const host=publicIP||location.hostname||'127.0.0.1'; const port=scheme==='http'?(stripColon((configCache&&configCache.http_port)||'7802')):(stripColon((configCache&&configCache.socks5_port)||'7801')); const url=scheme+'://'+encodeProxyUserInfo(user)+':'+encodeProxyUserInfo(pass)+'@'+host+':'+port; const okMsg=rawPass?'已复制':'已复制，请将 PASSWORD 替换为真实密码'; navigator.clipboard.writeText(url).then(()=>showToast(okMsg)).catch(()=>showToast('复制失败')) }
// toggleStar: 加星直接生效；取消星标须 confirm() 确认。
async function toggleStar(id,on){ if(on){ if(!confirm('取消该节点星标？'))return } return runAsync('星标操作失败',async()=>{ await api('/api/proxy/star',{method:'POST',body:JSON.stringify({id,starred:!on})}); await loadProxies(); showToast(on?'已取消星标':'已加星标') }) }
function renderProxies(){const protocol=document.getElementById('protocol-filter').value;const region=document.getElementById('region-filter').value;const sf=document.getElementById('status-filter').value;const srcf=(document.getElementById('source-filter')||{}).value||'';const qf=filterVal('quality-filter');const cff=filterVal('cf-filter');const aif={openai:filterVal('ai-openai-filter'),claude:filterVal('ai-claude-filter'),grok:filterVal('ai-grok-filter'),gemini:filterVal('ai-gemini-filter')};const latMinRaw=filterVal('latency-min');const latMaxRaw=filterVal('latency-max');const latMin=latMinRaw===''?null:Number(latMinRaw);const latMax=latMaxRaw===''?null:Number(latMaxRaw);const kw=filterVal('keyword-filter').toLowerCase();let rows=allProxies.filter(p=>(!protocol||p.protocol===protocol)&&(!region||regionOf(p)===region));if(sf)rows=rows.filter(p=>nodeState(p)===sf);if(srcf==='manual')rows=rows.filter(p=>p.source==='manual');else if(srcf==='subscription')rows=rows.filter(p=>p.source!=='manual');if(qf)rows=rows.filter(p=>qualityOf(p)===qf);if(cff)rows=rows.filter(p=>cfStateOf(p)===cff);['openai','claude','grok','gemini'].forEach(function(svc){const v=aif[svc];if(v)rows=rows.filter(p=>aiStateOf(p,svc)===v)});if(latMin!==null&&Number.isFinite(latMin))rows=rows.filter(p=>Number(p.latency||0)>=latMin);if(latMax!==null&&Number.isFinite(latMax))rows=rows.filter(p=>Number(p.latency||0)<=latMax);if(kw)rows=rows.filter(p=>{const addr=String(p.address||'').toLowerCase();const note=String(p.note||'').toLowerCase();return addr.indexOf(kw)>=0||note.indexOf(kw)>=0});const order={ok:0,pending:1,paused:2,failed:3};rows.sort((a,b)=>{const fa=(nodeState(a)==='ok'&&(a.starred===true||Number(a.starred)===1))?1:0;const fb=(nodeState(b)==='ok'&&(b.starred===true||Number(b.starred)===1))?1:0;if(fa!==fb)return fb-fa;const sa=nodeState(a),sb=nodeState(b);if(order[sa]!==order[sb])return order[sa]-order[sb];return Number(a.latency||1e9)-Number(b.latency||1e9)});const body=document.getElementById('proxy-rows');if(rows.length===0){body.innerHTML='<tr><td colspan="14" class="empty">没有匹配节点</td></tr>';return}proxyRenderRows=rows;proxyRenderCount=0;renderProxyBatch()}
function renderRegions(){const counts={};allProxies.filter(p=>isAvailable(p)&&isKnownRegion(p)).forEach(p=>{const r=regionOf(p);counts[r]=(counts[r]||0)+1});const entries=Object.keys(counts).sort().map(region=>({region,count:counts[region]}));const total=entries.reduce((sum,item)=>sum+item.count,0);document.getElementById('region-total').textContent=total+' 个可用节点';const list=document.getElementById('region-list');if(entries.length===0){list.innerHTML='<div class="empty">暂无可用地域数据</div>';return}list.innerHTML=entries.map(item=>{const pct=total?Math.round(item.count*100/total):0;return '<div class="region-row"><strong>'+html(item.region).toUpperCase()+'</strong><div class="bar"><span style="width:'+pct+'%"></span></div><span class="cnt">'+html(item.count)+'</span></div>'}).join('')}
// 总览节点分布：按地域+延迟档聚合圆点，有 session 的地域画连线。
// renderWorldMap 保留为旧调用名，实际转调 renderOrbitSystem。
const ORBIT_TRACKS={s:{rr:0.42,w:15,dir:1,phase:0},a:{rr:0.60,w:11,dir:-1,phase:40},b:{rr:0.78,w:8.5,dir:1,phase:15},c:{rr:0.96,w:6.5,dir:-1,phase:70}};
const ORBIT_QVAR={s:'var(--q-s)',a:'var(--q-a)',b:'var(--q-b)',c:'var(--q-c)'};
let orbitSats=[];let orbitT=0;let orbitLast=0;let orbitPaused=false;let orbitRAF=0;let orbitBuilt=false;
function orbitQualityTrack(p){const g=qualityOf(p);if(g==='S'||g==='A'||g==='B'||g==='C')return g.toLowerCase();const lat=Number(p&&p.latency||0);if(lat>0&&lat<=500)return 's';if(lat>0&&lat<=1000)return 'a';if(lat>0&&lat<=2000)return 'b';return 'c'}
function orbitStageGeom(){const st=document.getElementById('orbit-stage');const w=st?st.clientWidth:600;const h=st?st.clientHeight:338;return {cx:w/2,cy:h/2,halfW:w/2,halfH:h/2}}
function orbitAngAbsDiff(a,b){let d=a-b;while(d>Math.PI)d-=Math.PI*2;while(d<-Math.PI)d+=Math.PI*2;return Math.abs(d)}
function orbitRibbonPath(sx,sy,x,y,baseW,phase,widthScale,wind,lens){const dx=x-sx,dy=y-sy;const len=Math.hypot(dx,dy)||1;const ux=dx/len,uy=dy/len;const nx=-uy,ny=ux;const SEG=20;const swing=Math.min(len*0.038,6.5)*(0.85+0.15*(0.5+0.5*Math.sin(phase*0.32)));const side=Math.sin(phase*0.34);const wScale=widthScale==null?1:widthScale;const top=[],bot=[];let windHitMax=0;for(let i=0;i<=SEG;i++){const tt=i/SEG;let bend=swing*side*Math.sin(tt*Math.PI);let px=sx+ux*len*tt+nx*bend;let py=sy+uy*len*tt+ny*bend;let thin=1;if(wind&&wind.r>0){const rdx=px-wind.ox,rdy=py-wind.oy;const dist=Math.hypot(rdx,rdy)||1;const ang=Math.atan2(rdy,rdx);if(orbitAngAbsDiff(ang,wind.angle)<=wind.halfAperture){const u=(dist-wind.r)/Math.max(8,wind.band);const axis=1-orbitAngAbsDiff(ang,wind.angle)/Math.max(1e-4,wind.halfAperture);const hit=Math.exp(-u*u)*Math.pow(Math.max(0,axis),0.7);if(hit>0.02){if(hit>windHitMax)windHitMax=hit;const push=wind.force*hit;const tx=-wind.wy,ty=wind.wx;px+=wind.wx*push+tx*push*0.18*side;py+=wind.wy*push+ty*push*0.18*side;thin*=1-0.82*hit}}}if(lens&&lens.rx>0&&lens.ry>0){const ldx=px-lens.lx,ldy=py-lens.ly;const nr=Math.hypot(ldx/lens.rx,ldy/lens.ry);if(nr<1.35){const fall=Math.exp(-2.2*nr*nr);const w=fall*lens.strength;const rlen=Math.hypot(ldx,ldy)||1;const radial=Math.sin(lens.phase*0.7)*0.55;px+=(ldx/rlen)*w*radial*lens.rx*0.22;py+=(ldy/rlen)*w*radial*lens.ry*0.22;const txx=-ldy/rlen,tyy=ldx/rlen;const swirl=Math.sin(nr*Math.PI*1.1+lens.phase*0.9)*0.35;px+=txx*w*swirl*lens.rx*0.12;py+=tyy*w*swirl*lens.ry*0.12;thin*=1+0.12*fall*lens.strength}}const envelope=Math.pow(Math.sin(tt*Math.PI),0.95);const travel=0.92+0.08*Math.sin(tt*Math.PI-phase);const breath=0.97+0.03*Math.sin(phase*0.45);const hw=Math.max(0.12,(baseW*0.5)*envelope*travel*breath*wScale*Math.max(0.06,thin));top.push([px+nx*hw,py+ny*hw]);bot.push([px-nx*hw,py-ny*hw])}return {d:orbitRibbonSmooth(top,bot),windHit:windHitMax}}
function orbitRibbonSmooth(top,bot){function append(d,pts,move){if(pts.length<2)return d;if(move)d+='M'+pts[0][0].toFixed(1)+' '+pts[0][1].toFixed(1);for(let i=0;i<pts.length-1;i++){const p0=pts[Math.max(0,i-1)],p1=pts[i],p2=pts[i+1],p3=pts[Math.min(pts.length-1,i+2)];const c1x=p1[0]+(p2[0]-p0[0])/6,c1y=p1[1]+(p2[1]-p0[1])/6;const c2x=p2[0]-(p3[0]-p1[0])/6,c2y=p2[1]-(p3[1]-p1[1])/6;d+=' C'+c1x.toFixed(1)+' '+c1y.toFixed(1)+' '+c2x.toFixed(1)+' '+c2y.toFixed(1)+' '+p2[0].toFixed(1)+' '+p2[1].toFixed(1)}return d}let d=append('',top,true);d=append(d,bot.slice().reverse(),false);return d+' Z'}
function orbitSetGrad(id,sx,sy,x,y){const g=document.getElementById(id);if(!g)return;g.setAttribute('x1',sx.toFixed(1));g.setAttribute('y1',sy.toFixed(1));g.setAttribute('x2',x.toFixed(1));g.setAttribute('y2',y.toFixed(1))}
// 偶发装饰粒子流（约 5 分钟一次，首波约 18 秒）；纯视觉，不表示业务状态。
const SOLAR_WIND={active:false,front:0,duration:4.8,nextIn:18,period:300,strength:1,angle:0,halfAperture:0.18,band:38,streams:null};
function spawnWindStreams(){const n=7+Math.floor(Math.random()*5);const kinds=['spiral','hook','wave','s'];const arr=[];for(let i=0;i<n;i++){const t=(i/(n-1||1))-0.5;arr.push({da:t*SOLAR_WIND.halfAperture*2.1+(Math.random()-0.5)*0.04,w:0.9+Math.random()*1.6,op:0.22+Math.random()*0.38,core:Math.abs(t)<0.1,kind:kinds[Math.floor(Math.random()*kinds.length)],twist:(0.9+1.4*Math.random())*(Math.random()<0.5?-1:1),hook:0.55+0.9*Math.random(),waves:1.2+1.6*Math.random(),amp:0.10+0.18*Math.random(),phase:Math.random()*Math.PI*2,seed:Math.random()*10,curveSide:Math.random()<0.5?-1:1})}return arr}
function windPlumePath(cx,cy,halfW,halfH,angle,halfA,p){const edge=0.1,steps=16;const a0=angle-halfA*1.25,a1=angle+halfA*1.25;const x0=cx+halfW*edge*Math.cos(a0),y0=cy+halfH*edge*Math.sin(a0);const x3=cx+halfW*edge*Math.cos(a1),y3=cy+halfH*edge*Math.sin(a1);let d='M'+x0.toFixed(1)+' '+y0.toFixed(1);const n0x=-Math.sin(a0),n0y=Math.cos(a0);const bulge=Math.min(halfW,halfH)*0.05*p;const m0=0.45;const cx0=cx+halfW*(edge+(p-edge)*m0)*Math.cos(a0)+n0x*bulge;const cy0=cy+halfH*(edge+(p-edge)*m0)*Math.sin(a0)+n0y*bulge;const x1=cx+halfW*p*Math.cos(a0),y1=cy+halfH*p*Math.sin(a0);d+=' Q'+cx0.toFixed(1)+' '+cy0.toFixed(1)+' '+x1.toFixed(1)+' '+y1.toFixed(1);for(let i=1;i<=steps;i++){const tt=a0+(a1-a0)*i/steps;const rp=p*(1+0.03*Math.sin(i*1.7+angle));d+=' L'+(cx+halfW*rp*Math.cos(tt)).toFixed(1)+' '+(cy+halfH*rp*Math.sin(tt)).toFixed(1)}const n1x=-Math.sin(a1),n1y=Math.cos(a1);const cx1=cx+halfW*(edge+(p-edge)*m0)*Math.cos(a1)+n1x*bulge;const cy1=cy+halfH*(edge+(p-edge)*m0)*Math.sin(a1)+n1y*bulge;d+=' Q'+cx1.toFixed(1)+' '+cy1.toFixed(1)+' '+x3.toFixed(1)+' '+y3.toFixed(1)+' Z';return d}
function windStreamCurve(cx,cy,halfW,halfH,angle,p,s,timeP){const a0=angle+s.da;const edge=0.1;const SEG=10;const pts=[];const baseAmp=s.amp*Math.min(halfW,halfH)*Math.max(0.15,p);for(let i=0;i<=SEG;i++){const tt=i/SEG;const r=edge+(p-edge)*tt;let a=a0;let nOff=0;if(s.kind==='spiral'){a=a0+s.twist*tt*tt;nOff=baseAmp*(0.35+0.65*tt)*Math.sin(s.phase+tt*Math.PI*s.waves+timeP*1.8)}else if(s.kind==='hook'){const hookT=Math.max(0,(tt-0.55)/0.45);const hookEase=hookT*hookT*(3-2*hookT);a=a0+(s.curveSide||1)*s.hook*hookEase*1.1;nOff=baseAmp*0.45*Math.sin(s.phase+tt*2)*(1-hookEase*0.3)}else if(s.kind==='s'){a=a0+s.twist*0.25*Math.sin(tt*Math.PI);nOff=baseAmp*Math.sin(tt*Math.PI*2+s.phase+timeP)*(0.5+0.5*tt)}else{a=a0+s.da*0.2*Math.sin(tt*Math.PI);nOff=baseAmp*Math.sin(tt*Math.PI*s.waves+s.phase+timeP*2.1)*(0.4+0.6*tt)}nOff+=baseAmp*0.12*Math.sin(s.seed+tt*5+timeP*3);const ux=Math.cos(a),uy=Math.sin(a);const nx=-uy,ny=ux;pts.push([cx+halfW*r*ux+nx*nOff,cy+halfH*r*uy+ny*nOff])}if(pts.length<2)return '';let d='M'+pts[0][0].toFixed(1)+' '+pts[0][1].toFixed(1);for(let i=0;i<pts.length-1;i++){const p0=pts[Math.max(0,i-1)],p1=pts[i],p2=pts[i+1],p3=pts[Math.min(pts.length-1,i+2)];const c1x=p1[0]+(p2[0]-p0[0])/6,c1y=p1[1]+(p2[1]-p0[1])/6;const c2x=p2[0]-(p3[0]-p1[0])/6,c2y=p2[1]-(p3[1]-p1[1])/6;d+=' C'+c1x.toFixed(1)+' '+c1y.toFixed(1)+' '+c2x.toFixed(1)+' '+c2y.toFixed(1)+' '+p2[0].toFixed(1)+' '+p2[1].toFixed(1)}return d}
function updateSolarWind(dt,halfW,halfH,cx,cy){if(!SOLAR_WIND.active){SOLAR_WIND.nextIn-=dt;if(SOLAR_WIND.nextIn<=0){SOLAR_WIND.active=true;SOLAR_WIND.front=0;SOLAR_WIND.strength=0.85+0.15*Math.random();SOLAR_WIND.angle=Math.random()*Math.PI*2;SOLAR_WIND.halfAperture=(10+Math.random()*8)*Math.PI/180;SOLAR_WIND.band=32+18*Math.random();SOLAR_WIND.streams=spawnWindStreams();SOLAR_WIND.nextIn=SOLAR_WIND.period+(Math.random()*80-40)}}else{SOLAR_WIND.front+=dt/SOLAR_WIND.duration;if(SOLAR_WIND.front>=1){SOLAR_WIND.active=false;SOLAR_WIND.front=0;SOLAR_WIND.streams=null}}const g=document.getElementById('orbit-wind');const streamsG=document.getElementById('orbit-wind-streams');const plume=document.getElementById('orbit-wind-plume');if(g&&streamsG&&plume){if(SOLAR_WIND.active&&SOLAR_WIND.streams){const p=SOLAR_WIND.front;const ease=1-Math.pow(1-Math.min(1,p),1.45);const fade=Math.sin(Math.min(1,p)*Math.PI);const op=(0.4+0.55*fade)*SOLAR_WIND.strength;g.setAttribute('opacity',op.toFixed(3));const maxRx=halfW*1.12,maxRy=halfH*1.12;plume.setAttribute('d',windPlumePath(cx,cy,maxRx,maxRy,SOLAR_WIND.angle,SOLAR_WIND.halfAperture,ease));let streamHtml='';SOLAR_WIND.streams.forEach(s=>{if(s.curveSide==null)s.curveSide=s.twist>=0?1:-1;const d=windStreamCurve(cx,cy,maxRx,maxRy,SOLAR_WIND.angle,ease,s,p);if(!d)return;const sw=(1.1+s.w)*(0.75+0.25*fade);streamHtml+='<path class="orbit-wind-stream" d="'+d+'" stroke="url(#orbitWindStream)" stroke-width="'+(sw*2.2).toFixed(2)+'" stroke-opacity="'+(s.op*0.3*fade).toFixed(2)+'"/>';streamHtml+='<path class="'+(s.core?'orbit-wind-stream-core':'orbit-wind-stream')+'" d="'+d+'" stroke="'+(s.core?'#d8eaff':'url(#orbitWindStream)')+'" stroke-width="'+sw.toFixed(2)+'" stroke-opacity="'+(s.op*fade).toFixed(2)+'"/>'});streamsG.innerHTML=streamHtml;SOLAR_WIND._r=ease*Math.sqrt(halfW*halfH)*1.08;SOLAR_WIND._ox=cx;SOLAR_WIND._oy=cy;SOLAR_WIND._wx=Math.cos(SOLAR_WIND.angle);SOLAR_WIND._wy=Math.sin(SOLAR_WIND.angle)}else{g.setAttribute('opacity','0');streamsG.innerHTML='';plume.setAttribute('d','');SOLAR_WIND._r=0}}if(!SOLAR_WIND.active)return null;return{ox:SOLAR_WIND._ox,oy:SOLAR_WIND._oy,r:SOLAR_WIND._r||0,band:SOLAR_WIND.band,force:28*SOLAR_WIND.strength,angle:SOLAR_WIND.angle,halfAperture:SOLAR_WIND.halfAperture,wx:SOLAR_WIND._wx||0,wy:SOLAR_WIND._wy||0}}
// 偶发光晕扭曲（约 30 分钟一次，首波约 45 秒）；纯视觉，不表示业务状态。
const GRAV_LENS={active:false,life:0,duration:8,nextIn:45,period:1800,strength:0,phase:0,lx:0,ly:0,rx:0,ry:0};
function updateGravLens(dt,halfW,halfH,cx,cy){if(!GRAV_LENS.active){GRAV_LENS.nextIn-=dt;if(GRAV_LENS.nextIn<=0){GRAV_LENS.active=true;GRAV_LENS.life=0;GRAV_LENS.phase=Math.random()*Math.PI*2;const a=Math.random()*Math.PI*2;const rr=0.28+0.38*Math.random();GRAV_LENS.lx=cx+halfW*rr*Math.cos(a);GRAV_LENS.ly=cy+halfH*rr*Math.sin(a);GRAV_LENS.rx=halfW*(0.18+0.12*Math.random());GRAV_LENS.ry=halfH*(0.18+0.12*Math.random());GRAV_LENS.strength=0.85+0.15*Math.random();GRAV_LENS.nextIn=GRAV_LENS.period+(Math.random()*240-120)}}else{GRAV_LENS.life+=dt;GRAV_LENS.phase+=dt*0.85;if(GRAV_LENS.life>=GRAV_LENS.duration){GRAV_LENS.active=false;GRAV_LENS.life=0}}const g=document.getElementById('orbit-lens');const halo=document.getElementById('orbit-lens-halo');const rim=document.getElementById('orbit-lens-rim');if(g&&halo&&rim){if(GRAV_LENS.active){const p=GRAV_LENS.life/GRAV_LENS.duration;const env=Math.sin(Math.min(1,p)*Math.PI);const op=0.15+0.55*env*GRAV_LENS.strength;g.setAttribute('opacity',op.toFixed(3));const breathe=1+0.06*Math.sin(GRAV_LENS.phase*0.6);halo.setAttribute('cx',GRAV_LENS.lx.toFixed(1));halo.setAttribute('cy',GRAV_LENS.ly.toFixed(1));halo.setAttribute('rx',(GRAV_LENS.rx*breathe).toFixed(1));halo.setAttribute('ry',(GRAV_LENS.ry*breathe).toFixed(1));rim.setAttribute('cx',GRAV_LENS.lx.toFixed(1));rim.setAttribute('cy',GRAV_LENS.ly.toFixed(1));rim.setAttribute('rx',(GRAV_LENS.rx*breathe*0.92).toFixed(1));rim.setAttribute('ry',(GRAV_LENS.ry*breathe*0.92).toFixed(1))}else{g.setAttribute('opacity','0')}}if(!GRAV_LENS.active)return null;const p=GRAV_LENS.life/GRAV_LENS.duration;const env=Math.sin(Math.min(1,p)*Math.PI);return{lx:GRAV_LENS.lx,ly:GRAV_LENS.ly,rx:GRAV_LENS.rx,ry:GRAV_LENS.ry,strength:GRAV_LENS.strength*env,phase:GRAV_LENS.phase}}
function buildOrbitSvg(){const svg=document.getElementById('orbit-svg');if(!svg)return;const {cx,cy,halfW,halfH}=orbitStageGeom();let defs='<defs>';['s','a','b','c'].forEach(q=>{const c=getComputedStyle(document.documentElement).getPropertyValue('--q-'+q).trim()||'#3b8dff';const energy=getComputedStyle(document.documentElement).getPropertyValue('--accent').trim()||c;defs+='<linearGradient id="orbitBeam-'+q+'" gradientUnits="userSpaceOnUse" x1="0" y1="0" x2="1" y2="0"><stop offset="0%" stop-color="'+energy+'" stop-opacity="0"/><stop offset="12%" stop-color="'+energy+'" stop-opacity="0.75"/><stop offset="55%" stop-color="'+c+'" stop-opacity="0.85"/><stop offset="88%" stop-color="'+c+'" stop-opacity="0.55"/><stop offset="100%" stop-color="'+c+'" stop-opacity="0"/></linearGradient><linearGradient id="orbitGlow-'+q+'" gradientUnits="userSpaceOnUse" x1="0" y1="0" x2="1" y2="0"><stop offset="0%" stop-color="'+energy+'" stop-opacity="0"/><stop offset="20%" stop-color="'+c+'" stop-opacity="0.35"/><stop offset="70%" stop-color="'+c+'" stop-opacity="0.22"/><stop offset="100%" stop-color="'+c+'" stop-opacity="0"/></linearGradient>'});defs+='<linearGradient id="orbitWindStream" gradientUnits="userSpaceOnUse" x1="0" y1="0" x2="1" y2="0"><stop offset="0%" stop-color="#9ccaff" stop-opacity="0"/><stop offset="12%" stop-color="#b8d8ff" stop-opacity="0.55"/><stop offset="55%" stop-color="#6aa8ff" stop-opacity="0.28"/><stop offset="100%" stop-color="#3b8dff" stop-opacity="0"/></linearGradient>';defs+='<radialGradient id="orbitWindPlume" cx="50%" cy="50%" r="50%"><stop offset="0%" stop-color="#9ccaff" stop-opacity="0.2"/><stop offset="45%" stop-color="#5a9dff" stop-opacity="0.08"/><stop offset="100%" stop-color="#3b8dff" stop-opacity="0"/></radialGradient>';defs+='<radialGradient id="orbitLensFill" cx="50%" cy="50%" r="50%"><stop offset="0%" stop-color="#c8dcff" stop-opacity="0.18"/><stop offset="55%" stop-color="#8eb6ff" stop-opacity="0.08"/><stop offset="100%" stop-color="#3b8dff" stop-opacity="0"/></radialGradient>';defs+='</defs>';let rings='';['s','a','b','c'].forEach(q=>{const tr=ORBIT_TRACKS[q];const rx=halfW*tr.rr,ry=halfH*tr.rr;const c=getComputedStyle(document.documentElement).getPropertyValue('--q-'+q).trim()||'#3b8dff';rings+='<ellipse class="orbit-ring" cx="'+cx+'" cy="'+cy+'" rx="'+rx.toFixed(1)+'" ry="'+ry.toFixed(1)+'" stroke="'+c+'" stroke-opacity="0.34"/>'});const wind='<g id="orbit-wind" opacity="0"><path class="orbit-wind-plume" id="orbit-wind-plume" d=""/><g id="orbit-wind-streams"></g></g>';const lens='<g id="orbit-lens" opacity="0"><ellipse class="orbit-lens-halo" id="orbit-lens-halo" cx="'+cx+'" cy="'+cy+'" rx="0" ry="0"/><ellipse class="orbit-lens-rim" id="orbit-lens-rim" cx="'+cx+'" cy="'+cy+'" rx="0" ry="0"/></g>';svg.innerHTML=defs+rings+wind+lens+'<g id="orbit-beams"></g>'}
function buildOrbitSats(){const layer=document.getElementById('orbit-sats');const beamG=document.getElementById('orbit-beams');if(!layer||!beamG)return;layer.innerHTML='';beamG.innerHTML='';orbitSats=[];const sessCount={};(Array.isArray(orbitSessions)?orbitSessions:[]).forEach(s=>{const r=String((s&&s.region)||'').trim().toLowerCase();if(!r||r==='unknown')return;sessCount[r]=(sessCount[r]||0)+1});const buckets={};allProxies.filter(p=>isAvailable(p)&&isKnownRegion(p)).forEach(p=>{const cc=regionOf(p);const q=orbitQualityTrack(p);const key=cc+'|'+q;if(!buckets[key])buckets[key]={cc:cc,q:q,n:0,k:0};buckets[key].n++});Object.keys(buckets).forEach(key=>{const b=buckets[key];b.k=sessCount[b.cc]||0});const byQ={s:[],a:[],b:[],c:[]};Object.values(buckets).forEach(b=>{if(byQ[b.q])byQ[b.q].push(b)});const svgns='http://www.w3.org/2000/svg';['s','a','b','c'].forEach(q=>{const arr=byQ[q];if(!arr.length)return;arr.sort((x,y)=>y.n-x.n);const tr=ORBIT_TRACKS[q];const step=360/arr.length;arr.forEach((d,i)=>{const el=document.createElement('div');el.className='orbit-sat'+(d.k>0?' live':'');el.style.setProperty('--qc',ORBIT_QVAR[q]);const SMIN=30,SMAX=60,NLO=1,NHI=40;const norm=Math.max(0,Math.min(1,(Math.sqrt(d.n)-Math.sqrt(NLO))/(Math.sqrt(NHI)-Math.sqrt(NLO))));const size=Math.round(SMIN+(SMAX-SMIN)*norm);el.dataset.size=String(size);el.style.width=size+'px';el.style.height=size+'px';const tip=html(d.cc).toUpperCase()+' · '+html(d.n)+' 节点 · '+q.toUpperCase()+' 档'+(d.k>0?(' · '+html(d.k)+' 会话'):'');el.innerHTML='<div class="ball"><span class="cc">'+html(d.cc).toUpperCase()+'</span></div><span class="cnt num">'+html(d.n)+'</span>';el.title=tip;layer.appendChild(el);let beam=null;if(d.k>0){const g=document.createElementNS(svgns,'g');const glow=document.createElementNS(svgns,'path');glow.setAttribute('class','orbit-beam-glow');glow.setAttribute('fill','url(#orbitGlow-'+q+')');const path=document.createElementNS(svgns,'path');path.setAttribute('class','orbit-beam');path.setAttribute('fill','url(#orbitBeam-'+q+')');g.appendChild(glow);g.appendChild(path);beamG.appendChild(g);beam={path:path,glow:glow,phase:Math.random()*Math.PI*2,speed:1.1+0.2*Math.min(5,d.k),baseW:Math.max(2.2,Math.min(5.5,2.0+0.85*Math.min(6,d.k)))}}orbitSats.push({el:el,beam:beam,track:tr,baseAngle:tr.phase+i*step,q:q})})});const ipEl=document.getElementById('orbit-gw-ip');if(ipEl){ipEl.textContent=publicIP||(location&&location.hostname)||'--'}orbitBuilt=true}
function orbitFrame(now){if(!orbitLast)orbitLast=now;const dt=(now-orbitLast)/1000;orbitLast=now;const live=(!orbitPaused&&!document.hidden);if(live)orbitT+=dt;const {cx,cy,halfW,halfH}=orbitStageGeom();const sunR=42;const wind=live?updateSolarWind(dt,halfW,halfH,cx,cy):null;const lens=live?updateGravLens(dt,halfW,halfH,cx,cy):null;orbitSats.forEach(s=>{const tr=s.track;const ang=(s.baseAngle+tr.dir*tr.w*orbitT)*Math.PI/180;const rx=halfW*tr.rr,ry=halfH*tr.rr;const x=cx+rx*Math.cos(ang);const y=cy+ry*Math.sin(ang);const depth=(Math.sin(ang)+1)/2;const scale=0.82+0.30*depth;const size=Number(s.el.dataset.size)||28;s.el.style.left=x+'px';s.el.style.top=y+'px';s.el.style.width=(size*scale)+'px';s.el.style.height=(size*scale)+'px';s.el.style.zIndex=String(Math.round(depth*100)+10);if(s.beam){const dx=x-cx,dy=y-cy;const len=Math.hypot(dx,dy)||1;const sx=cx+dx/len*sunR,sy=cy+dy/len*sunR;s.beam.phase+=dt*s.beam.speed;const baseOp=0.42+0.48*depth;const main=orbitRibbonPath(sx,sy,x,y,s.beam.baseW,s.beam.phase,1,wind,lens);const hit=main.windHit||0;const op=baseOp*(1-0.78*hit);s.beam.path.setAttribute('d',main.d);s.beam.path.style.opacity=Math.max(0.08,op).toFixed(2);if(s.beam.glow){const g=orbitRibbonPath(sx,sy,x,y,s.beam.baseW,s.beam.phase,2.1,wind,lens);s.beam.glow.setAttribute('d',g.d);s.beam.glow.style.opacity=Math.max(0.04,op*0.5*(1-0.5*hit)).toFixed(2)}orbitSetGrad('orbitBeam-'+s.q,sx,sy,x,y);orbitSetGrad('orbitGlow-'+s.q,sx,sy,x,y)}});orbitRAF=requestAnimationFrame(orbitFrame)}
function ensureOrbitLoop(){if(!orbitRAF){orbitLast=0;orbitRAF=requestAnimationFrame(orbitFrame)}}
function renderOrbitSystem(){const stage=document.getElementById('orbit-stage');if(!stage)return;buildOrbitSvg();buildOrbitSats();ensureOrbitLoop()}
function toggleOrbitPause(){orbitPaused=!orbitPaused;const btn=document.getElementById('orbit-pause-btn');if(btn)btn.textContent=orbitPaused?'继续':'暂停'}
function renderWorldMap(){renderOrbitSystem()}

function expandAllSessions(open){document.querySelectorAll('#session-rows .session-card').forEach(function(el){el.classList.toggle('open',!!open)})}
function toggleSessionCard(el){if(!el)return;el.classList.toggle('open')}
async function loadSessions(){const sessions=await api('/api/sessions');if(!sessions)return;orbitSessions=Array.isArray(sessions)?sessions:[];renderOrbitSystem();const body=document.getElementById('session-rows');const cnt=document.getElementById('sess-count');if(cnt)cnt.textContent=Array.isArray(sessions)?(sessions.length+' 条 sticky 绑定'):'--';if(!Array.isArray(sessions)||sessions.length===0){body.innerHTML='<div class="empty">暂无活跃 session</div>';return}body.innerHTML=sessions.map(function(s){const sid=html(s.session_id);const masked=html(maskAddress(s.node));const region=String(s.region||'').trim().toLowerCase();const regionBadge=region&&region!=='unknown'?'<span class="badge ok">'+html(region).toUpperCase()+'</span> ':'<span class="badge gray">未知</span> ';const ttlSec=Number(s.remaining_ttl_seconds)||0;const ttlCls=ttlSec>0&&ttlSec<60?' danger':(ttlSec>0&&ttlSec<180?' warn':'');return '<div class="session-card"><div class="head" onclick="toggleSessionCard(this.parentElement)"><span class="sid" title="'+sid+'">'+sid+'</span><div class="chips">'+regionBadge+'</div><span class="ttl'+ttlCls+'">'+html(formatTTL(ttlSec))+'</span><span class="chev" aria-hidden="true">›</span></div><div class="body"><div class="detail-grid"><div class="di"><span class="k">会话 ID</span><span class="v mono">'+sid+'</span></div><div class="di"><span class="k">出口地域</span><span class="v">'+regionBadge+'</span></div><div class="di"><span class="k">出口节点</span><span class="v mono">'+regionBadge+masked+'</span></div><div class="di"><span class="k">剩余 TTL</span><span class="v">'+html(formatTTL(ttlSec))+'</span></div></div></div></div>'}).join('')}
function formatTTL(seconds){const value=Number(seconds)||0;const min=Math.floor(value/60);const sec=value%60;return min>0?min+'m '+sec+'s':sec+'s'}
async function addManualNode(){return runAsync('添加失败',async()=>{const payload={link:document.getElementById('manual-link').value.trim(),region:document.getElementById('manual-region').value.trim(),note:document.getElementById('manual-note').value.trim()};if(!payload.link){showToast('请填写节点链接');return}await api('/api/manual-node/add',{method:'POST',body:JSON.stringify(payload)});document.getElementById('manual-link').value='';document.getElementById('manual-region').value='';document.getElementById('manual-note').value='';await Promise.all([loadStats(),loadProxies()]);showToast('手工节点已添加')})}
function toggleSelectAll(on){document.querySelectorAll('.proxy-select').forEach(el=>{el.checked=!!on})}
function selectedProxyIDs(){return Array.from(document.querySelectorAll('.proxy-select:checked')).map(el=>Number(el.value)).filter(n=>Number.isFinite(n)&&n>0)}
async function batchDeleteSelected(){return runAsync('批量删除失败',async()=>{const ids=selectedProxyIDs();if(!ids.length){showToast('请先勾选手工节点');return}if(!confirm('删除选中的 '+ids.length+' 个手工节点？'))return;const r=await api('/api/manual-node/batch-delete',{method:'POST',body:JSON.stringify({ids})});await Promise.all([loadStats(),loadProxies()]);showToast('已删除 '+(r&&r.deleted!=null?r.deleted:ids.length)+' 个'+(r&&r.failed?('，失败 '+r.failed):''))})}
async function importManualNodes(){return runAsync('批量导入失败',async()=>{const text=document.getElementById('import-text').value;const region=document.getElementById('import-region').value.trim();const note=document.getElementById('import-note').value.trim();if(!String(text||'').trim()){showToast('请粘贴代理列表');return}const r=await api('/api/manual-node/import',{method:'POST',body:JSON.stringify({text,region,note})});document.getElementById('import-modal').classList.remove('show');document.getElementById('import-text').value='';await Promise.all([loadStats(),loadProxies()]);showToast('导入完成：新增 '+(r.added||0)+' / 跳过 '+(r.skipped||0)+' / 失败 '+(r.failed||0))})}
async function manageManualNode(id,address){return runAsync('管理失败',async()=>{const current=allProxies.find(p=>Number(p.id)===Number(id))||{};const choice=prompt('手工节点管理：输入 1=改地域，2=改备注，3=删除', '1');if(choice===null)return;if(choice==='1'){const region=prompt('地域',current.region||'');if(region===null)return;await api('/api/manual-node/region',{method:'POST',body:JSON.stringify({id,address,region})});await loadProxies();showToast('地域已更新');return}if(choice==='2'){const note=prompt('备注',current.note||'');if(note===null)return;await api('/api/manual-node/note',{method:'POST',body:JSON.stringify({id,address,note})});await loadProxies();showToast('备注已更新');return}if(choice==='3'){if(!confirm('删除此手工节点？'))return;await api('/api/manual-node/delete',{method:'POST',body:JSON.stringify({id,address})});await Promise.all([loadStats(),loadProxies()]);showToast('手工节点已删除')}})}
async function editManualRegion(id,address){return runAsync('地域更新失败',async()=>{const current=allProxies.find(p=>Number(p.id)===Number(id))||{};const region=prompt('地域',current.region||'');if(region===null)return;await api('/api/manual-node/region',{method:'POST',body:JSON.stringify({id,address,region})});await loadProxies();showToast('地域已更新')})}
async function editManualNote(id,address){return runAsync('备注更新失败',async()=>{const current=allProxies.find(p=>Number(p.id)===Number(id))||{};const note=prompt('备注',current.note||'');if(note===null)return;await api('/api/manual-node/note',{method:'POST',body:JSON.stringify({id,address,note})});await loadProxies();showToast('备注已更新')})}
async function deleteManualNode(id,address){return runAsync('删除失败',async()=>{if(!confirm('删除此手工节点？'))return;await api('/api/manual-node/delete',{method:'POST',body:JSON.stringify({id,address})});await Promise.all([loadStats(),loadProxies()]);showToast('手工节点已删除')})}
async function toggleProxy(id,address,enable){return runAsync('操作失败',async()=>{await api('/api/proxy/toggle',{method:'POST',body:JSON.stringify({id,address,enable})});await Promise.all([loadStats(),loadProxies()]);showToast(enable?'节点已启用':'节点已停用')})}
// testProxy: 触发单节点重新验证（走完整 ValidateOne，含连通 google/openai/github/cloudflare/gstatic），后端异步执行，稍后自动刷新列表。
async function testProxy(id,address){return runAsync('测试失败',async()=>{await api('/api/proxy/refresh',{method:'POST',body:JSON.stringify({id,address})});showToast('测试连通已启动，稍后自动刷新');setTimeout(()=>runAsync('刷新失败',()=>Promise.all([loadStats(),loadProxies()])),4000)})}
async function loadSubscriptions(){const subs=await api('/api/subscriptions');if(!subs)return;const box=document.getElementById('sub-list');if(!Array.isArray(subs)||subs.length===0){box.innerHTML='<div class="empty">暂无订阅，点右上角“添加订阅”</div>';return}box.innerHTML=subs.map(sub=>{const paused=sub.status==='paused';const activeCount=Number(sub.active_count||0);const disabledCount=Number(sub.disabled_count||0);const proxyCount=Number(sub.proxy_count||0);const pausedCount=Number(sub.paused_count??Math.max(0,proxyCount-activeCount-disabledCount));const toggleLabel=paused?'启用':'暂停';const badge=paused?'<span class="badge warn">已暂停</span>':'<span class="badge ok">活跃</span>';const id=Number(sub.id);const idArg=Number.isFinite(id)?String(id):'0';return '<div class="sub-item"><div class="meta"><strong>'+html(sub.name)+' '+badge+'</strong><div class="muted">'+html(activeCount)+' 可用 / '+html(pausedCount)+' 暂停 / '+html(disabledCount)+' 不可用</div></div><div class="mini-actions"><button class="mini" onclick="refreshSub('+idArg+')" title="重新拉取并验证">刷新</button><button class="mini" onclick="toggleSub('+idArg+')" title="启用或暂停该订阅及其节点">'+toggleLabel+'</button><button class="mini danger" onclick="deleteSub('+idArg+')" title="删除订阅及其节点">删除</button></div></div>'}).join('')}
function openSubModal(){document.getElementById('sub-modal').classList.add('show')}function closeSubModal(){document.getElementById('sub-modal').classList.remove('show')}
async function addSubscription(){return runAsync('添加失败',async()=>{const payload={name:document.getElementById('sub-name').value.trim(),url:document.getElementById('sub-url').value.trim(),file_content:document.getElementById('sub-file-content').value.trim(),headers:document.getElementById('sub-headers').value.trim(),refresh_min:Number(document.getElementById('sub-refresh').value)||60};if(!payload.url&&!payload.file_content){showToast('请填写订阅 URL 或粘贴配置内容');return}await api('/api/subscription/add',{method:'POST',body:JSON.stringify(payload)});document.getElementById('sub-name').value='';document.getElementById('sub-url').value='';document.getElementById('sub-file-content').value='';document.getElementById('sub-headers').value='';closeSubModal();await Promise.all([loadSubscriptions(),loadStats(),loadProxies()]);showToast('订阅已添加')})}
async function refreshSub(id){return runAsync('刷新失败',async()=>{await api('/api/subscription/refresh',{method:'POST',body:JSON.stringify({id})});showToast('刷新已启动，稍后自动更新');refreshLater()})}
async function refreshAllSubs(){return runAsync('刷新失败',async()=>{await api('/api/subscription/refresh-all',{method:'POST'});showToast('全部刷新已启动，稍后自动更新');refreshLater()})}
async function toggleSub(id){return runAsync('切换失败',async()=>{await api('/api/subscription/toggle',{method:'POST',body:JSON.stringify({id})});await Promise.all([loadSubscriptions(),loadStats(),loadProxies()]);showToast('已切换启用/暂停状态')})}
async function deleteSub(id){return runAsync('删除失败',async()=>{if(!confirm('删除此订阅及其全部节点？'))return;await api('/api/subscription/delete',{method:'POST',body:JSON.stringify({id})});await Promise.all([loadSubscriptions(),loadStats(),loadProxies()]);showToast('订阅已删除')})}
async function loadLogs(){const data=await api('/api/logs');if(!data)return;const box=document.getElementById('logs-box');if(!box)return;const prevTop=box.scrollTop;const lines=Array.isArray(data.lines)?data.lines:[];box.innerHTML=lines.length?lines.map(line=>'<div class="log-line">'+html(line)+'</div>').join(''):'<div class="log-line">no logs</div>';const auto=document.getElementById('logs-autoscroll');if(auto&&auto.checked){box.scrollTop=box.scrollHeight}else{box.scrollTop=prevTop}}
async function loadConfig(){configCache=await api('/api/config');if(!configCache)return;const hp=stripColon(configCache.http_port),sp=stripColon(configCache.socks5_port),wp=stripColon(configCache.webui_port);document.getElementById('cfg-http-port').value=hp;document.getElementById('cfg-socks5-port').value=sp;document.getElementById('cfg-webui-port').value=wp;document.getElementById('cfg-auth-enabled').value=String(Boolean(configCache.proxy_auth_enabled));document.getElementById('cfg-auth-username').value=configCache.proxy_auth_username||'';document.getElementById('cfg-auth-password').value='';document.getElementById('cfg-session-ttl').value=configCache.session_ttl_minutes||'';document.getElementById('cfg-default-region').value=configCache.default_region||'';document.getElementById('cfg-health-interval').value=configCache.health_check_interval||'';document.getElementById('cfg-max-retry').value=configCache.max_retry??'';document.getElementById('cfg-singbox-path').value=configCache.singbox_path||'';document.getElementById('cfg-allowed-countries').value=(configCache.allowed_countries||[]).join(',');document.getElementById('cfg-blocked-countries').value=(configCache.blocked_countries||[]).join(',');renderConnection();renderDSLExamples()}
async function loadPublicIP(){return runAsync('公网 IP 获取失败',async()=>{const d=await api('/api/public-ip');if(d){if(d.public_ip){publicIP=d.public_ip;renderConnection()}if(d.country){gatewayCC=String(d.country).toLowerCase()}renderOrbitSystem()}})}
function renderConnection(){if(!configCache)return;const sp=stripColon(configCache.socks5_port)||'7801';const hp=stripColon(configCache.http_port)||'7802';const base=configCache.proxy_auth_username||'username';const enabled=configCache.proxy_auth_enabled;const host=publicIP||location.hostname||'127.0.0.1';document.getElementById('conn-socks5').textContent=host+':'+(sp||'7801');document.getElementById('conn-http').textContent=host+':'+(hp||'7802');document.getElementById('conn-user').textContent=base;document.getElementById('conn-pass').textContent=enabled?'见首次启动日志 / 系统设置':'（认证已关闭，无需密码）';document.getElementById('conn-auth-state').textContent=enabled?'代理认证：开启':'代理认证：关闭';const cred=enabled?(base+':PASSWORD@'):'';document.getElementById('conn-cmd').textContent='curl --socks5 '+cred+host+':'+(sp||'7801')+' https://www.gstatic.com/generate_204'}
function renderDSLExamples(){const base=(configCache&&configCache.proxy_auth_username)?configCache.proxy_auth_username:'username';const box=document.getElementById('dsl-examples');if(box){box.innerHTML=['-region-us','-unlock-gpt','-region-jp-unlock-all-session-app01','-session-browser'].map(s=>'<div class="guide-row"><b>'+html(base)+'</b><span>'+html(s)+'</span></div>').join('')}const hint=document.getElementById('dsl-hint');if(hint){hint.textContent=(configCache&&configCache.proxy_auth_enabled)?('前缀 “'+base+'” = 代理认证用户名；-region-XX 地域；-unlock-gpt|claude|gemini|grok|cf|all 解锁过滤；-session-ID 黏连。'):'代理认证当前关闭；启用后前缀须等于代理认证用户名。'}}
async function openSettings(){switchTab('settings')}function closeSettings(){}function countries(id){return document.getElementById(id).value.split(',').map(v=>v.trim().toUpperCase()).filter(Boolean)}
function formatAPIKeyTime(v){if(!v)return '--';const d=new Date(v);return Number.isNaN(d.getTime())?String(v):d.toLocaleString()}
function renderAPIKeys(keys){const body=document.getElementById('apikey-rows');if(!body)return;const list=Array.isArray(keys)?keys:[];if(!list.length){body.innerHTML='<tr><td colspan="5" class="empty">暂无 API Key</td></tr>';return}body.innerHTML=list.map(k=>{const id=html(k.id);const name=html(k.name);const created=html(formatAPIKeyTime(k.created_at));const last=html(formatAPIKeyTime(k.last_used_at));const disabled=!!(k.disabled===true||Number(k.disabled)===1);const st=disabled?'<span class="badge warn">已吊销</span>':'<span class="badge ok">有效</span>';const revokeBtn=disabled?'':'<button class="mini" onclick="revokeAPIKey(\''+id+'\')">吊销</button> ';return '<tr><td>'+name+'</td><td>'+created+'</td><td>'+last+'</td><td>'+st+'</td><td>'+revokeBtn+'<button class="mini danger" onclick="deleteAPIKey(\''+id+'\')">删除</button></td></tr>'}).join('')}
async function loadAPIKeys(){const data=await api('/api/apikeys');if(!data)return;renderAPIKeys(data.keys||data||[])}
async function createAPIKey(){return runAsync('创建 API Key 失败',async()=>{const name=document.getElementById('apikey-name').value.trim();if(!name){showToast('请填写 Key 名称');return}const r=await api('/api/apikey/create',{method:'POST',body:JSON.stringify({name})});document.getElementById('apikey-name').value='';document.getElementById('apikey-once-name').value=r&&r.name?r.name:name;document.getElementById('apikey-once-key').value=r&&r.key?r.key:'';document.getElementById('apikey-once-modal').classList.add('show');await loadAPIKeys();showToast('API Key 已创建（仅显示一次）')})}
async function revokeAPIKey(id){return runAsync('吊销失败',async()=>{if(!confirm('吊销该 API Key？吊销后立即失效。'))return;await api('/api/apikey/revoke',{method:'POST',body:JSON.stringify({id})});await loadAPIKeys();showToast('已吊销')})}
async function deleteAPIKey(id){return runAsync('删除失败',async()=>{if(!confirm('删除该 API Key？此操作不可恢复。'))return;await api('/api/apikey/delete',{method:'POST',body:JSON.stringify({id})});await loadAPIKeys();showToast('已删除')})}
async function saveConfig(){return runAsync('保存失败',async()=>{if(!configCache)await loadConfig();if(!configCache)throw new Error('配置未加载');const payload={proxy_auth_enabled:document.getElementById('cfg-auth-enabled').value==='true',proxy_auth_username:document.getElementById('cfg-auth-username').value.trim(),proxy_auth_password:document.getElementById('cfg-auth-password').value,session_ttl_minutes:Number(document.getElementById('cfg-session-ttl').value),default_region:document.getElementById('cfg-default-region').value.trim().toLowerCase(),health_check_interval:Number(document.getElementById('cfg-health-interval').value),max_retry:Number(document.getElementById('cfg-max-retry').value),singbox_path:document.getElementById('cfg-singbox-path').value.trim(),allowed_countries:countries('cfg-allowed-countries'),blocked_countries:countries('cfg-blocked-countries')};await api('/api/config/save',{method:'POST',body:JSON.stringify(payload)});await loadConfig();showToast('配置已保存')})}
// ===== 侧边栏折叠持久化 =====
function applySidebar(collapsed){document.body.classList.toggle('sidebar-collapsed',!!collapsed);try{localStorage.setItem('gg-sidebar',collapsed?'1':'0')}catch(e){}}
function toggleSidebar(){applySidebar(!document.body.classList.contains('sidebar-collapsed'))}
function openDrawer(){document.body.classList.add('drawer-open')}
function closeDrawer(){document.body.classList.remove('drawer-open')}
(function(){let c=false;try{c=localStorage.getItem('gg-sidebar')==='1'}catch(e){}applySidebar(c);const sb=document.getElementById('sidebar');if(sb)requestAnimationFrame(function(){sb.classList.remove('preload')})})();
// AI/Cloudflare 图标筛选：点击循环 全部->畅通->阻断->未知；值写入隐藏 select，renderProxies 读取不变。
const FILTER_CYCLE={'':'全部','unlocked':'畅通','blocked':'阻断','unprobed':'未知','unknown':'未知'};
function cycleFilter(selId,btnId){const sel=document.getElementById(selId);if(!sel)return;const opts=Array.from(sel.options).map(o=>o.value);let idx=opts.indexOf(sel.value);idx=(idx+1)%opts.length;sel.value=opts[idx];syncFilterToggle(selId,btnId);renderProxies()}
function syncFilterToggle(selId,btnId){const sel=document.getElementById(selId);const btn=document.getElementById(btnId);if(!sel||!btn)return;const v=sel.value;const st=btn.querySelector('.st');if(st)st.textContent=FILTER_CYCLE[v]||'全部';btn.setAttribute('aria-pressed',v?'true':'false')}
function initFilterToggles(){document.querySelectorAll('.filter-toggle[data-sel]').forEach(function(btn){syncFilterToggle(btn.dataset.sel,btn.id)})}
// 节点表分批渲染：首批立即渲染，滚动接近底部再增量，避免上千行一次性 DOM。
let proxyRenderRows=[];let proxyRenderCount=0;const PROXY_BATCH=80;
function proxyRowHTML(p){const addr=addressArg(p.address);const id=proxyIDArg(p);const manual=p.source==='manual';const st=nodeState(p);const label=html(nodeLabel(p));const showRegion=isAvailable(p)&&isKnownRegion(p);const toggleBtn=(st==='paused')?'<button class="mini" onclick="toggleProxy('+id+',decodeURIComponent(\''+addr+'\'),true)">启用</button>':'<button class="mini" onclick="toggleProxy('+id+',decodeURIComponent(\''+addr+'\'),false)">停用</button>';const testBtn='<button class="mini" onclick="testProxy('+id+',decodeURIComponent(\''+addr+'\'))">测试</button>';const copyBtn='<button class="mini" onclick="copyProxyCred('+id+')">复制</button>';const baseActions=testBtn+' '+copyBtn+' '+toggleBtn;const manageBtn=manual?('<button class="mini" onclick="manageManualNode('+id+',decodeURIComponent(\''+addr+'\'))">管理</button>'):'';const actions=baseActions+(manageBtn?(' '+manageBtn):'');const latencyText=Number(p.latency)>0?html(p.latency)+' ms':'--';const sel=manual?'<input type="checkbox" class="proxy-select" value="'+id+'">':'';return '<tr><td>'+sel+'</td><td>'+starBtn(p)+'</td><td title="'+label+'">'+label+'</td><td>'+protocolBadges(p)+'</td><td>'+(showRegion?'<span class="badge ok">'+html(regionOf(p)).toUpperCase()+'</span>':'<span class="muted">--</span>')+'</td><td class="mono">'+html(p.exit_ip)+'</td><td>'+latencyText+'</td><td>'+abuserBadge(p.ipapiis_score)+'</td><td>'+ipapiFlagsBadges(p.ipapi_flags,!!p.ipapi_flags_seen)+'</td><td>'+cfBadge(p.cf_blocked)+'</td><td>'+aiBadges(p.ai_reachability)+'</td><td>'+html(sourceLabel(p))+'</td><td>'+stateBadge(st)+'</td><td>'+actions+'</td></tr>'}
function renderProxyBatch(){const body=document.getElementById('proxy-rows');if(!body)return;const next=Math.min(proxyRenderCount+PROXY_BATCH,proxyRenderRows.length);let h='';for(let i=proxyRenderCount;i<next;i++)h+=proxyRowHTML(proxyRenderRows[i]);if(proxyRenderCount===0)body.innerHTML=h;else body.insertAdjacentHTML('beforeend',h);proxyRenderCount=next}
function onProxyScroll(){if(proxyRenderCount>=proxyRenderRows.length)return;const el=document.documentElement;if(el.scrollTop+el.clientHeight>=el.scrollHeight-320)renderProxyBatch()}
window.addEventListener('scroll',onProxyScroll,{passive:true});
let orbitResizeTimer=null;window.addEventListener('resize',function(){clearTimeout(orbitResizeTimer);orbitResizeTimer=setTimeout(function(){try{renderOrbitSystem()}catch(e){}},180)});
// 骨架墓碑：载入态灰条 shimmer（尊重 prefers-reduced-motion，动画由 CSS 关闭）。
function skeletonRows(n){let h='';for(let i=0;i<(n||3);i++)h+='<div class="skeleton sk-row"></div>';return '<div class="skeleton-wrap">'+h+'</div>'}
function showSkeletons(){['region-list','sub-list','session-rows','singbox-status'].forEach(function(id){const el=document.getElementById(id);if(el)el.innerHTML=skeletonRows(3)})}
// 首次进入总览再画分布图。
let overviewSeen=false;function markViewLazy(name){if(name==='overview'&&!overviewSeen){overviewSeen=true;try{renderOrbitSystem()}catch(e){}}}
initFilterToggles();showSkeletons();markViewLazy('overview');
refreshAll();
setInterval(()=>runAsync('自动刷新失败',()=>Promise.all([loadStats(),loadProxies(),loadSubscriptions(),loadSessions()])),10000);
setInterval(()=>runAsync('日志刷新失败',loadLogs),5000);`
