-- 直播礼物初始数据 (来源: docs/gift.txt)
-- 表: live_gifts
-- name=中文名称, name_en=英文, name_es=西班牙语, name_pt=葡萄牙语, name_hi=印地语
-- sort=36..1 (sort desc 展示，低价礼物 sort 更大，价格升序在前)

DELETE FROM `live_gifts` WHERE `name_en` IN (
'magis stick', 'party cake', 'lemon Tea', 'chocolate', 'Ice Cream', 'finger heart', 'lipstick', 'kiss', 'bubble shoot', 'loveU', 'birthday gift', 'blue rose', 'royal salute', 'love baloon', 'gown', 'lollipop', 'carousel', 'ring box', 'carriage', 'ametyhst', 'heart of sea', 'moonlight', 'cupid''s bow', 'ferris wheel', 'sea castle', 'shooting star', 'fantasy city', 'lamborghini', 'helicopter', 'summer yacth', 'rolls royce', 'air plane', 'air balloon', 'luxury yacth', 'star ship', 'rocket'
);

INSERT INTO `live_gifts` (`name`, `name_en`, `name_es`, `name_pt`, `name_hi`, `icon`, `animation`, `price`, `category`, `sort`, `status`, `published_at`, `description`, `created_at`, `updated_at`) VALUES
('魔法棒', 'magis stick', 'varita mágica', 'varinha mágica', 'जादुई छड़ी', '', '', 10, '', 36, 1, NOW(), '', NOW(), NOW()),
('派对蛋糕', 'party cake', 'pastel de fiesta', 'bolo de festa', 'पार्टी केक', '', '', 20, '', 35, 1, NOW(), '', NOW(), NOW()),
('柠檬茶', 'lemon Tea', 'té de limón', 'chá de limão', 'नींबू चाय', '', '', 40, '', 34, 1, NOW(), '', NOW(), NOW()),
('黑巧克力', 'chocolate', 'chocolate negro', 'chocolate amargo', 'डार्क चॉकलेट', '', '', 50, '', 33, 1, NOW(), '', NOW(), NOW()),
('冰淇凌', 'Ice Cream', 'helado', 'sorvete', 'आइसक्रीम', '', '', 60, '', 32, 1, NOW(), '', NOW(), NOW()),
('比心', 'finger heart', 'corazón con los dedos', 'coração com os dedos', 'उंगली से दिल बनाना', '', '', 70, '', 31, 1, NOW(), '', NOW(), NOW()),
('口红', 'lipstick', 'labial', 'batom', 'लिपस्टिक', '', '', 80, '', 30, 1, NOW(), '', NOW(), NOW()),
('吻', 'kiss', 'beso', 'beijo', 'चुंबन', '', '', 100, '', 29, 1, NOW(), '', NOW(), NOW()),
('吹泡泡', 'bubble shoot', 'soplar burbujas', 'soprar bolhas', 'बुलबुले उड़ाना', '', '', 120, '', 28, 1, NOW(), '', NOW(), NOW()),
('一箭穿心', 'loveU', 'flecha atravesando el corazón', 'flecha atravessando o coração', 'दिल में तीर', '', '', 150, '', 27, 1, NOW(), '', NOW(), NOW()),
('生日祝福', 'birthday gift', 'regalo de cumpleaños', 'presente de aniversário', 'जन्मदिन उपहार', '', '', 170, '', 26, 1, NOW(), '', NOW(), NOW()),
('爱心蓝玫瑰', 'blue rose', 'rosa azul con corazón', 'rosa azul com coração', 'दिल वाला नीला गुलाब', '', '', 200, '', 25, 1, NOW(), '', NOW(), NOW()),
('礼炮', 'royal salute', 'salva real', 'salva real', 'शाही तोप का सलाम', '', '', 250, '', 24, 1, NOW(), '', NOW(), NOW()),
('爱心气球', 'love baloon', 'globo de corazón', 'balão de coração', 'दिल वाला गुब्बारा', '', '', 300, '', 23, 1, NOW(), '', NOW(), NOW()),
('婚纱', 'gown', 'vestido de novia', 'vestido de noiva', 'शादी का गाउन', '', '', 350, '', 22, 1, NOW(), '', NOW(), NOW()),
('棒棒糖', 'lollipop', 'paleta de caramelo', 'pirulito', 'लॉलीपॉप', '', '', 400, '', 21, 1, NOW(), '', NOW(), NOW()),
('旋转木马', 'carousel', 'carrusel', 'carrossel', 'घूमता हुआ घोड़ा', '', '', 500, '', 20, 1, NOW(), '', NOW(), NOW()),
('戒指盒子', 'ring box', 'caja de anillo', 'caixa de anel', 'अंगूठी का डिब्बा', '', '', 600, '', 19, 1, NOW(), '', NOW(), NOW()),
('马车', 'carriage', 'carruaje', 'carruagem', 'घोड़ा गाड़ी', '', '', 700, '', 18, 1, NOW(), '', NOW(), NOW()),
('紫水晶', 'ametyhst', 'amatista', 'ametista', 'जामुनिया रत्न', '', '', 850, '', 17, 1, NOW(), '', NOW(), NOW()),
('海洋之心', 'heart of sea', 'corazón del mar', 'coração do mar', 'समुद्र का दिल', '', '', 1000, '', 16, 1, NOW(), '', NOW(), NOW()),
('月光', 'moonlight', 'luz de luna', 'luz da lua', 'चंद्रमा की रोशनी', '', '', 1300, '', 15, 1, NOW(), '', NOW(), NOW()),
('丘比特之剑', 'cupid''s bow', 'flecha de cupido', 'flecha de cupido', 'कामदेव का तीर', '', '', 1500, '', 14, 1, NOW(), '', NOW(), NOW()),
('摩天轮', 'ferris wheel', 'noria', 'roda gigante', 'फेरिस व्हील', '', '', 2000, '', 13, 1, NOW(), '', NOW(), NOW()),
('海洋城堡', 'sea castle', 'castillo marítimo', 'castelo marinho', 'समुद्री महल', '', '', 2500, '', 12, 1, NOW(), '', NOW(), NOW()),
('流星雨', 'shooting star', 'lluvia de estrellas', 'chuva de estrelas', 'उल्का बारिश', '', '', 3200, '', 11, 1, NOW(), '', NOW(), NOW()),
('城堡', 'fantasy city', 'ciudad fantástica', 'cidade fantasia', 'काल्पनिक शहर', '', '', 4500, '', 10, 1, NOW(), '', NOW(), NOW()),
('兰博基尼', 'lamborghini', 'lamborghini', 'lamborghini', 'लैंबोर्गिनी', '', '', 6500, '', 9, 1, NOW(), '', NOW(), NOW()),
('直升机', 'helicopter', 'helicóptero', 'helicóptero', 'हेलीकॉप्टर', '', '', 9000, '', 8, 1, NOW(), '', NOW(), NOW()),
('夏季游艇', 'summer yacth', 'yate de verano', 'iate de verão', 'गर्मी की यॉट', '', '', 12500, '', 7, 1, NOW(), '', NOW(), NOW()),
('劳斯莱斯', 'rolls royce', 'rolls royce', 'rolls royce', 'रोल्स रॉयस', '', '', 17500, '', 6, 1, NOW(), '', NOW(), NOW()),
('飞机', 'air plane', 'avión', 'avião', 'हवाई जहाज', '', '', 25000, '', 5, 1, NOW(), '', NOW(), NOW()),
('热气球', 'air balloon', 'globo aerostático', 'balão de ar quente', 'गर्म हवा का गुब्बारा', '', '', 35000, '', 4, 1, NOW(), '', NOW(), NOW()),
('豪华邮轮', 'luxury yacth', 'yate de lujo', 'iate de luxo', 'लक्जरी यॉट', '', '', 50000, '', 3, 1, NOW(), '', NOW(), NOW()),
('宇宙飞船', 'star ship', 'nave estelar', 'nave estelar', 'तारों का जहाज', '', '', 75000, '', 2, 1, NOW(), '', NOW(), NOW()),
('火箭', 'rocket', 'cohete', 'foguete', 'रॉकेट', '', '', 120000, '', 1, 1, NOW(), '', NOW(), NOW());
