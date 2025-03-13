package data

// Company consists of company information
var Companys = map[string][]string{
	"name":         {"{person.en_US_last} {company.suffix}", "{person.en_US_last}-{person.en_US_last}", "{person.en_US_last}, {person.en_US_last} and {person.en_US_last}"},
	"suffix":       {"Inc", "and Sons", "LLC", "Group"},
	"buzzwords":    {"Adaptive", "Advanced", "Ameliorated", "Assimilated", "Automated", "Balanced", "Business-focused", "Centralized", "Cloned", "Compatible", "Configurable", "Cross-group", "Cross-platform", "Customer-focused", "Customizable", "De-engineered", "Decentralized", "Devolved", "Digitized", "Distributed", "Diverse", "Down-sized", "Enhanced", "Enterprise-wide", "Ergonomic", "Exclusive", "Expanded", "Extended", "Face to face", "Focused", "Front-line", "Fully-configurable", "Function-based", "Fundamental", "Future-proofed", "Grass-roots", "Horizontal", "Implemented", "Innovative", "Integrated", "Intuitive", "Inverse", "Managed", "Mandatory", "Monitored", "Multi-channelled", "Multi-lateral", "Multi-layered", "Multi-tiered", "Networked", "Object-based", "Open-architected", "Open-source", "Operative", "Optimized", "Optional", "Organic", "Organized", "Persevering", "Persistent", "Phased", "Polarised", "Pre-emptive", "Proactive", "Profit-focused", "Profound", "Programmable", "Progressive", "Public-key", "Quality-focused", "Re-contextualized", "Re-engineered", "Reactive", "Realigned", "Reduced", "Reverse-engineered", "Right-sized", "Robust", "Seamless", "Secured", "Self-enabling", "Sharable", "Stand-alone", "Streamlined", "Switchable", "Synchronised", "Synergistic", "Synergized", "Team-oriented", "Total", "Triple-buffered", "Universal", "Up-sized", "Upgradable", "User-centric", "User-friendly", "Versatile", "Virtual", "Vision-oriented", "Visionary", "24 hour", "24/7", "3rd generation", "4th generation", "5th generation", "6th generation", "actuating", "analyzing", "asymmetric", "asynchronous", "attitude-oriented", "background", "bandwidth-monitored", "bi-directional", "bifurcated", "bottom-line", "clear-thinking", "client-driven", "client-server", "coherent", "cohesive", "composite", "content-based", "context-sensitive", "contextually-based", "dedicated", "demand-driven", "didactic", "directional", "discrete", "disintermediate", "dynamic", "eco-centric", "empowering", "encompassing", "even-keeled", "executive", "explicit", "exuding", "fault-tolerant", "foreground", "fresh-thinking", "full-range", "global", "grid-enabled", "heuristic", "high-level", "holistic", "homogeneous", "human-resource", "hybrid", "impactful", "incremental", "intangible", "interactive", "intermediate", "leading edge", "local", "logistical", "maximized", "methodical", "mission-critical", "mobile", "modular", "motivating", "multi-state", "multi-tasking", "multimedia", "national", "needs-based", "neutral", "next generation", "non-volatile", "object-oriented", "optimal", "optimizing", "radical", "real-time", "reciprocal", "regional", "responsive", "scalable", "secondary", "solution-oriented", "stable", "static", "system-worthy", "systematic", "systemic", "tangible", "tertiary", "transitional", "uniform", "upward-trending", "user-facing", "value-added", "web-enabled", "well-modulated", "zero administration", "zero defect", "zero tolerance", "Graphic Interface", "Graphical User Interface", "ability", "access", "adapter", "algorithm", "alliance", "analyzer", "application", "approach", "architecture", "archive", "array", "artificial intelligence", "attitude", "benchmark", "budgetary management", "capability", "capacity", "challenge", "circuit", "collaboration", "complexity", "concept", "conglomeration", "contingency", "core", "customer loyalty", "data-warehouse", "database", "definition", "emulation", "encoding", "encryption", "extranet", "firmware", "flexibility", "focus group", "forecast", "frame", "framework", "function", "functionalities", "groupware", "hardware", "help-desk", "hierarchy", "hub", "implementation", "info-mediaries", "infrastructure", "initiative", "installation", "instruction set", "interface", "internet solution", "intranet", "knowledge base", "knowledge user", "leverage", "local area network", "matrices", "matrix", "methodology", "middleware", "migration", "model", "moderator", "monitoring", "moratorium", "neural-net", "open architecture", "open system", "orchestration", "paradigm", "parallelism", "policy", "portal", "pricing structure", "process improvement", "product", "productivity", "project", "projection", "protocol", "secured line", "service-desk", "software", "solution", "standardization", "strategy", "structure", "success", "superstructure", "support", "synergy", "system engine", "task-force", "throughput", "time-frame", "toolset", "utilisation", "website", "workforce"},
	"bs":           {"aggregate", "architect", "benchmark", "brand", "cultivate", "deliver", "deploy", "disintermediate", "drive", "e-enable", "embrace", "empower", "enable", "engage", "engineer", "enhance", "envisioneer", "evolve", "expedite", "exploit", "extend", "facilitate", "generate", "grow", "harness", "implement", "incentivize", "incubate", "innovate", "integrate", "iterate", "leverage", "matrix", "maximize", "mesh", "monetize", "morph", "optimize", "orchestrate", "productize", "recontextualize", "redefine", "reintermediate", "reinvent", "repurpose", "revolutionize", "scale", "seize", "strategize", "streamline", "syndicate", "synergize", "synthesize", "target", "transform", "transition", "unleash", "utilize", "visualize", "whiteboard", "24/365", "24/7", "B2B", "B2C", "back-end", "best-of-breed", "bleeding-edge", "bricks-and-clicks", "clicks-and-mortar", "collaborative", "compelling", "cross-media", "cross-platform", "customized", "cutting-edge", "distributed", "dot-com", "dynamic", "e-business", "efficient", "end-to-end", "enterprise", "extensible", "frictionless", "front-end", "global", "granular", "holistic", "impactful", "innovative", "integrated", "interactive", "intuitive", "killer", "leading-edge", "magnetic", "mission-critical", "next-generation", "one-to-one", "open-source", "out-of-the-box", "plug-and-play", "proactive", "real-time", "revolutionary", "rich", "robust", "scalable", "seamless", "sexy", "sticky", "strategic", "synergistic", "transparent", "turn-key", "ubiquitous", "user-centric", "value-added", "vertical", "viral", "virtual", "visionary", "web-enabled", "wireless", "world-class", "ROI", "action-items", "applications", "architectures", "bandwidth", "channels", "communities", "content", "convergence", "deliverables", "e-business", "e-commerce", "e-markets", "e-services", "e-tailers", "experiences", "eyeballs", "functionalities", "infomediaries", "infrastructures", "initiatives", "interfaces", "markets", "methodologies", "metrics", "mindshare", "models", "networks", "niches", "paradigms", "partnerships", "platforms", "portals", "relationships", "schemas", "solutions", "supply-chains", "synergies", "systems", "technologies", "users", "vortals", "web services", "web-readiness"},
	"zh_CN":        {"{company.zh_CN_prefix} {company.zh_CN_suffix}"},
	"zh_CN_prefix": {"超艺", "和泰", "九方", "鑫博腾飞", "戴硕电子", "济南亿次元", "海创", "创联世纪", "凌云", "泰麒麟", "彩虹", "兰金电子", "晖来计算机", "天益", "恒聪百汇", "菊风公司", "惠派国际公司", "创汇", "思优", "时空盒数字", "易动力", "飞海科技", "华泰通安", "盟新", "商软冠联", "图龙信息", "易动力", "华远软件", "创亿", "时刻", "开发区世创", "明腾", "良诺", "天开", "毕博诚", "快讯", "凌颖信息", "黄石金承", "恩悌", "雨林木风计算机", "双敏电子", "维旺明", "网新恒天", "数字100", "飞利信", "立信电子", "联通时科", "中建创业", "新格林耐特", "新宇龙信息", "浙大万朋", "MBP软件", "昂歌信息", "万迅电脑", "方正科技", "联软", "七喜", "南康", "银嘉", "巨奥", "佳禾", "国讯", "信诚致远", "浦华众城", "迪摩", "太极", "群英", "合联电子", "同兴万点", "襄樊地球村", "精芯", "艾提科信", "昊嘉", "鸿睿思博", "四通", "富罳", "商软冠联", "诺依曼软件", "东方峻景", "华成育卓", "趋势", "维涛", "通际名联"},
	"zh_CN_suffix": {"科技有限公司", "网络有限公司", "信息有限公司", "传媒有限公司"},
}
