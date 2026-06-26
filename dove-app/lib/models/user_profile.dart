const kMaxProfilePhotos = 9;

class ProfilePrompt {
  const ProfilePrompt({
    required this.title,
    required this.content,
  });

  final String title;
  final String content;
}

class ProfileInfo {
  const ProfileInfo({
    required this.prompts,
    required this.details,
  });

  final List<ProfilePrompt> prompts;
  final List<(String label, String value)> details;
}

const mockProfileInfo = ProfileInfo(
  prompts: [
    ProfilePrompt(
      title: '关于我',
      content: '喜欢独立书店、爵士乐和周末徒步。相信好的关系从真诚的对话开始。',
    ),
    ProfilePrompt(
      title: '我在找',
      content: '能一起探索城市、也享受安静午后的人。',
    ),
    ProfilePrompt(
      title: '我的日常',
      content: '工作日做产品，晚上跑步或听黑胶。周末常去徐汇滨江，或者找一家没去过的小馆子。',
    ),
    ProfilePrompt(
      title: '最近迷上',
      content: '手冲咖啡、城市骑行，还有收集各城市的独立杂志。',
    ),
    ProfilePrompt(
      title: '周末理想状态',
      content: '睡到自然醒，买一束花，约朋友看展，晚上在家做饭。',
    ),
    ProfilePrompt(
      title: '一次难忘的经历',
      content: '在京都雨夜误入一家爵士 bar，和陌生人聊到关门，至今还记得那杯 highball 的味道。',
    ),
  ],
  details: [
    ('年龄', '28'),
    ('城市', '上海'),
    ('职业', '产品设计师'),
    ('身高', '178 cm'),
    ('学历', '本科'),
    ('家乡', '杭州'),
    ('吸烟', '不吸烟'),
    ('饮酒', '偶尔'),
  ],
);
