const APP_CONFIG = {
  appName: 'Sara Live',
  supportEmail: 'support@saralive.app',
  privacyEmail: 'privacy@saralive.app',
  apiBaseUrl: '',
};

const SITE_I18N = {
  en: {
    navHome: 'Home',
    navAbout: 'About',
    navFeatures: 'Features',
    navPrivacy: 'Privacy',
    navTerms: 'Terms',
    navCreatorTerms: 'Creator Terms',
    navRoomOwnerTerms: 'Room Owner',
    navVipDesc: 'VIP',
    navCancel: 'Delete Account',
    heroTitle: 'Watch and Go Live Anytime',
    heroDesc: 'Sara Live is a mobile live streaming platform for hosts and audiences. Enjoy real-time video, interactive chat, virtual gifts, and short video content in one app.',
    badgeLive: 'HD Live Streaming',
    badgeChat: 'Real-time Chat',
    badgeGift: 'Virtual Gifts',
    featuresTitle: 'Built for Live Entertainment',
    featuresDesc: 'Everything you need to discover live rooms, interact with creators, and share moments.',
    feat1Title: 'Live Broadcast',
    feat1Desc: 'Start or join live streams with low-latency video and audio.',
    feat2Title: 'Interactive Community',
    feat2Desc: 'Chat, follow creators, send gifts, and join the audience in real time.',
    feat3Title: 'Privacy & Safety',
    feat3Desc: 'We use moderation tools and clear policies to protect users and content.',
    footerRights: '© 2026 Sara Live. All rights reserved.',
    privacyPageTitle: 'Privacy Policy',
    privacyPageDesc: 'This Privacy Policy explains how Sara Live collects, uses, stores, shares and protects personal information, in compliance with global requirements including GDPR, LGPD, Middle East data regulations and Indonesian PSE rules.',
    legalTocTitle: 'Contents',
    legalTocToggle: 'Table of contents',
    legalContactLabel: 'Contact',
    termsPageTitle: 'Terms of Service',
    termsPageDesc: 'Rules and conditions for using Sara Live live streaming and related features.',
    creatorTermsPageTitle: 'Short Video Creator Upload Compliance Terms',
    creatorTermsPageDesc: 'Special terms for Sara Live short video creators and MCN-affiliated creators regarding content upload compliance, commercial disclosure, regional rules, and intellectual property.',
    roomOwnerTermsPageTitle: 'Room Owner Responsibility Terms',
    roomOwnerTermsPageDesc: 'Responsibilities and compliance requirements for Sara Live voice chat room and live room owners.',
    vipDescPageTitle: 'VIP Membership Description',
    vipDescPageDesc: 'Learn about Sara Live VIP levels, upgrade rules, and the privileges available at each tier.',
    lastUpdated: 'Last updated: July 27, 2026',
    cancelPageTitle: 'Delete Account',
    cancelPageDesc: 'Enter your deactivation code to permanently delete your Sara Live account. This action cannot be undone.',
    cancelCodeLabel: 'Deactivation code',
    cancelCodePlaceholder: 'Enter your deactivation code',
    cancelCodeHint: 'Find your deactivation code in the Sara Live app under account settings.',
    cancelSubmit: 'Delete my account',
    cancelConfirm: 'Are you sure you want to permanently delete this account? This cannot be undone.',
    cancelSuccess: 'Your account has been deleted successfully.',
    cancelFail: 'Unable to delete account. Please check your deactivation code and try again.',
    cancelBlocked: 'Too many failed attempts. Please try again in 2 hours.',
    cancelAlready: 'This account has already been deleted.',
    cancelNetworkError: 'Network error. Please try again later.',
    aboutPageTitle: 'About Sara Live',
    aboutP1: 'Sara Live is a live-streaming and social platform intended only for users aged 18 and older. It enables adults around the world to connect, communicate, and share experiences in real time.',
    aboutP2: 'Sara Live offers live video, voice rooms, calls, public chat, direct messaging, short videos, creator tools, and virtual gifting. Feature availability may vary by country, account status, and app version.',
    aboutP3: 'We are committed to maintaining a respectful and responsible community. Users can report content or accounts, block other users, and review our Community Guidelines and Child Safety Standards in the Safety Center.',
    aboutSupportTitle: 'Support and contact',
    aboutSupportDesc: 'Sara Live is operated by the developer identified in the applicable app-store listing. For product support, safety concerns, or legal inquiries, use Help & Support in Sara Live. If the in-app channel is unavailable, use the developer contact shown in the app-store listing from which you downloaded Sara Live.',
    aboutThanks: 'Thank you for being part of the Sara Live community.',
    aboutSupportLink: 'Support',
  },
  es: {
    navHome: 'Inicio',
    navAbout: 'Acerca de',
    navFeatures: 'Funciones',
    navPrivacy: 'Privacidad',
    navTerms: 'Términos',
    navCreatorTerms: 'Creadores',
    navRoomOwnerTerms: 'Propietario',
    navVipDesc: 'VIP',
    navCancel: 'Eliminar cuenta',
    heroTitle: 'Mira y transmite en vivo en cualquier momento',
    heroDesc: 'Sara Live es una plataforma móvil de transmisión en vivo para anfitriones y audiencias. Disfruta video en tiempo real, chat interactivo, regalos virtuales y videos cortos.',
    badgeLive: 'Transmisión HD',
    badgeChat: 'Chat en tiempo real',
    badgeGift: 'Regalos virtuales',
    featuresTitle: 'Diseñada para entretenimiento en vivo',
    featuresDesc: 'Todo lo que necesitas para descubrir salas en vivo, interactuar con creadores y compartir momentos.',
    feat1Title: 'Transmisión en vivo',
    feat1Desc: 'Inicia o únete a transmisiones con video y audio de baja latencia.',
    feat2Title: 'Comunidad interactiva',
    feat2Desc: 'Chatea, sigue creadores, envía regalos y participa en tiempo real.',
    feat3Title: 'Privacidad y seguridad',
    feat3Desc: 'Usamos moderación y políticas claras para proteger usuarios y contenido.',
    footerRights: '© 2026 Sara Live. Todos los derechos reservados.',
    privacyPageTitle: 'Política de Privacidad',
    privacyPageDesc: 'Esta Política de Privacidad explica cómo Sara Live recopila, utiliza, almacena, comparte y protege la información personal, cumpliendo requisitos globales como GDPR, LGPD, normativas de datos de Oriente Medio y cumplimiento PSE de Indonesia.',
    legalTocTitle: 'Contenido',
    legalTocToggle: 'Tabla de contenido',
    legalContactLabel: 'Contacto',
    termsPageTitle: 'Términos de Servicio',
    termsPageDesc: 'Reglas y condiciones para usar Sara Live y sus funciones relacionadas.',
    creatorTermsPageTitle: 'Términos de cumplimiento para subida de videos cortos',
    creatorTermsPageDesc: 'Términos especiales para creadores de videos cortos de Sara Live y creadores afiliados a MCN sobre cumplimiento de contenido, divulgación comercial, normas regionales y propiedad intelectual.',
    roomOwnerTermsPageTitle: 'Términos de responsabilidad del propietario de sala',
    roomOwnerTermsPageDesc: 'Responsabilidades y requisitos de cumplimiento para propietarios de salas de chat de voz y salas en vivo de Sara Live.',
    vipDescPageTitle: 'Descripción de membresía VIP',
    vipDescPageDesc: 'Conoce los niveles VIP de Sara Live, las reglas de mejora y los privilegios disponibles en cada nivel.',
    lastUpdated: 'Última actualización: 27 de julio de 2026',
    cancelPageTitle: 'Eliminar cuenta',
    cancelPageDesc: 'Introduce tu código de desactivación para eliminar permanentemente tu cuenta de Sara Live. Esta acción no se puede deshacer.',
    cancelCodeLabel: 'Código de desactivación',
    cancelCodePlaceholder: 'Introduce tu código de desactivación',
    cancelCodeHint: 'Encuentra tu código de desactivación en la app Sara Live, en ajustes de cuenta.',
    cancelSubmit: 'Eliminar mi cuenta',
    cancelConfirm: '¿Seguro que deseas eliminar permanentemente esta cuenta? Esta acción no se puede deshacer.',
    cancelSuccess: 'Tu cuenta se eliminó correctamente.',
    cancelFail: 'No se pudo eliminar la cuenta. Verifica el código e inténtalo de nuevo.',
    cancelBlocked: 'Demasiados intentos fallidos. Inténtalo de nuevo en 2 horas.',
    cancelAlready: 'Esta cuenta ya fue eliminada.',
    cancelNetworkError: 'Error de red. Inténtalo más tarde.',
    aboutPageTitle: 'Acerca de Sara Live',
    aboutP1: 'Sara Live es una plataforma de transmisión en vivo y social destinada únicamente a usuarios mayores de 18 años. Permite a adultos de todo el mundo conectarse, comunicarse y compartir experiencias en tiempo real.',
    aboutP2: 'Sara Live ofrece video en vivo, salas de voz, llamadas, chat público, mensajes directos, videos cortos, herramientas para creadores y regalos virtuales. La disponibilidad de funciones puede variar según el país, el estado de la cuenta y la versión de la app.',
    aboutP3: 'Estamos comprometidos con mantener una comunidad respetuosa y responsable. Los usuarios pueden denunciar contenido o cuentas, bloquear a otros usuarios y revisar nuestras Directrices de la Comunidad y Estándares de Seguridad Infantil en el Centro de Seguridad.',
    aboutSupportTitle: 'Soporte y contacto',
    aboutSupportDesc: 'Sara Live es operada por el desarrollador identificado en el listado correspondiente de la tienda de aplicaciones. Para soporte del producto, inquietudes de seguridad o consultas legales, usa Ayuda y Soporte en Sara Live. Si el canal dentro de la app no está disponible, utiliza el contacto del desarrollador que aparece en la tienda desde la que descargaste Sara Live.',
    aboutThanks: 'Gracias por ser parte de la comunidad Sara Live.',
    aboutSupportLink: 'Soporte',
  },
  pt: {
    navHome: 'Início',
    navAbout: 'Sobre',
    navFeatures: 'Recursos',
    navPrivacy: 'Privacidade',
    navTerms: 'Termos',
    navCreatorTerms: 'Criadores',
    navRoomOwnerTerms: 'Proprietário',
    navVipDesc: 'VIP',
    navCancel: 'Excluir conta',
    heroTitle: 'Assista e transmita ao vivo a qualquer hora',
    heroDesc: 'A Sara Live é uma plataforma móvel de transmissão ao vivo para hosts e público. Aproveite vídeo em tempo real, chat interativo, presentes virtuais e vídeos curtos.',
    badgeLive: 'Transmissão HD',
    badgeChat: 'Chat em tempo real',
    badgeGift: 'Presentes virtuais',
    featuresTitle: 'Feita para entretenimento ao vivo',
    featuresDesc: 'Tudo o que você precisa para descobrir salas ao vivo, interagir com criadores e compartilhar momentos.',
    feat1Title: 'Transmissão ao vivo',
    feat1Desc: 'Inicie ou entre em transmissões com vídeo e áudio de baixa latência.',
    feat2Title: 'Comunidade interativa',
    feat2Desc: 'Converse, siga criadores, envie presentes e participe em tempo real.',
    feat3Title: 'Privacidade e segurança',
    feat3Desc: 'Usamos moderação e políticas claras para proteger usuários e conteúdo.',
    footerRights: '© 2026 Sara Live. Todos os direitos reservados.',
    privacyPageTitle: 'Política de Privacidade',
    privacyPageDesc: 'Esta Política de Privacidade explica como a Sara Live coleta, usa, armazena, compartilha e protege informações pessoais, em conformidade com requisitos globais incluindo GDPR, LGPD, regulamentos de dados do Oriente Médio e conformidade PSE da Indonésia.',
    legalTocTitle: 'Conteúdo',
    legalTocToggle: 'Índice',
    legalContactLabel: 'Contato',
    termsPageTitle: 'Termos de Serviço',
    termsPageDesc: 'Regras e condições para usar a Sara Live e recursos relacionados.',
    creatorTermsPageTitle: 'Termos de conformidade para upload de vídeos curtos',
    creatorTermsPageDesc: 'Termos especiais para criadores de vídeos curtos da Sara Live e criadores afiliados a MCN sobre conformidade de conteúdo, divulgação comercial, regras regionais e propriedade intelectual.',
    roomOwnerTermsPageTitle: 'Termos de responsabilidade do proprietário de sala',
    roomOwnerTermsPageDesc: 'Responsabilidades e requisitos de conformidade para proprietários de salas de chat de voz e salas ao vivo da Sara Live.',
    vipDescPageTitle: 'Descrição da associação VIP',
    vipDescPageDesc: 'Saiba sobre os níveis VIP da Sara Live, regras de upgrade e privilégios disponíveis em cada nível.',
    lastUpdated: 'Última atualização: 27 de julho de 2026',
    cancelPageTitle: 'Excluir conta',
    cancelPageDesc: 'Digite seu código de desativação para excluir permanentemente sua conta Sara Live. Esta ação não pode ser desfeita.',
    cancelCodeLabel: 'Código de desativação',
    cancelCodePlaceholder: 'Digite seu código de desativação',
    cancelCodeHint: 'Encontre seu código de desativação no app Sara Live, em configurações da conta.',
    cancelSubmit: 'Excluir minha conta',
    cancelConfirm: 'Tem certeza de que deseja excluir permanentemente esta conta? Esta ação não pode ser desfeita.',
    cancelSuccess: 'Sua conta foi excluída com sucesso.',
    cancelFail: 'Não foi possível excluir a conta. Verifique o código e tente novamente.',
    cancelBlocked: 'Muitas tentativas falhadas. Tente novamente em 2 horas.',
    cancelAlready: 'Esta conta já foi excluída.',
    cancelNetworkError: 'Erro de rede. Tente novamente mais tarde.',
    aboutPageTitle: 'Sobre a Sara Live',
    aboutP1: 'A Sara Live é uma plataforma de transmissão ao vivo e social destinada apenas a usuários com 18 anos ou mais. Ela permite que adultos de todo o mundo se conectem, se comuniquem e compartilhem experiências em tempo real.',
    aboutP2: 'A Sara Live oferece vídeo ao vivo, salas de voz, chamadas, chat público, mensagens diretas, vídeos curtos, ferramentas para criadores e presentes virtuais. A disponibilidade de recursos pode variar conforme o país, o status da conta e a versão do app.',
    aboutP3: 'Estamos comprometidos em manter uma comunidade respeitosa e responsável. Os usuários podem denunciar conteúdo ou contas, bloquear outros usuários e revisar nossas Diretrizes da Comunidade e Padrões de Segurança Infantil no Centro de Segurança.',
    aboutSupportTitle: 'Suporte e contato',
    aboutSupportDesc: 'A Sara Live é operada pelo desenvolvedor identificado no respectivo listing da loja de aplicativos. Para suporte do produto, preocupações de segurança ou consultas legais, use Ajuda e Suporte na Sara Live. Se o canal no app não estiver disponível, use o contato do desenvolvedor exibido na loja de onde você baixou a Sara Live.',
    aboutThanks: 'Obrigado por fazer parte da comunidade Sara Live.',
    aboutSupportLink: 'Suporte',
  },
  hi: {
    navHome: 'होम',
    navAbout: 'परिचय',
    navFeatures: 'फीचर्स',
    navPrivacy: 'गोपनीयता',
    navTerms: 'नियम',
    navCreatorTerms: 'क्रिएटर',
    navRoomOwnerTerms: 'रूम मालिक',
    navVipDesc: 'VIP',
    navCancel: 'खाता हटाएँ',
    heroTitle: 'कभी भी देखें और लाइव जाएँ',
    heroDesc: 'Sara Live होस्ट और दर्शकों के लिए एक मोबाइल लाइव स्ट्रीमिंग प्लेटफ़ॉर्म है। रीयल-टाइम वीडियो, इंटरैक्टिव चैट, वर्चुअल गिफ्ट और शॉर्ट वीडियो का आनंद लें।',
    badgeLive: 'HD लाइव स्ट्रीमिंग',
    badgeChat: 'रीयल-टाइम चैट',
    badgeGift: 'वर्चुअल गिफ्ट',
    featuresTitle: 'लाइव मनोरंजन के लिए बनाया गया',
    featuresDesc: 'लाइव रूम खोजने, क्रिएटर्स के साथ जुड़ने और पल साझा करने के लिए सब कुछ।',
    feat1Title: 'लाइव प्रसारण',
    feat1Desc: 'कम विलंबता वाले वीडियो और ऑडियो के साथ लाइव स्ट्रीम शुरू या जॉइन करें।',
    feat2Title: 'इंटरैक्टिव समुदाय',
    feat2Desc: 'चैट करें, क्रिएटर्स को फॉलो करें, गिफ्ट भेजें और रीयल-टाइम में भाग लें।',
    feat3Title: 'गोपनीयता और सुरक्षा',
    feat3Desc: 'हम उपयोगकर्ताओं और सामग्री की सुरक्षा के लिए मॉडरेशन और स्पष्ट नीतियाँ उपयोग करते हैं।',
    footerRights: '© 2026 Sara Live. सर्वाधिकार सुरक्षित।',
    privacyPageTitle: 'गोपनीयता नीति',
    privacyPageDesc: 'यह गोपनीयता नीति बताती है कि Sara Live GDPR, LGPD, मध्य पूर्व डेटा विनियमों और इंडोनेशियाई PSE अनुपालन सहित वैश्विक आवश्यकताओं के अनुसार व्यक्तिगत जानकारी कैसे एकत्र, उपयोग, संग्रहीत, साझा और सुरक्षित करती है।',
    legalTocTitle: 'विषय सूची',
    legalTocToggle: 'विषय सूची',
    legalContactLabel: 'संपर्क',
    termsPageTitle: 'सेवा की शर्तें',
    termsPageDesc: 'Sara Live और संबंधित सुविधाओं के उपयोग के नियम और शर्तें।',
    creatorTermsPageTitle: 'शॉर्ट वीडियो क्रिएटर अपलोड अनुपालन शर्तें',
    creatorTermsPageDesc: 'Sara Live शॉर्ट वीडियो क्रिएटर्स और MCN-संबद्ध क्रिएटर्स के लिए सामग्री अनुपालन, वाणिज्यिक प्रकटीकरण, क्षेत्रीय नियम और बौद्धिक संपदा से संबंधित विशेष शर्तें।',
    roomOwnerTermsPageTitle: 'रूम मालिक जिम्मेदारी शर्तें',
    roomOwnerTermsPageDesc: 'Sara Live वॉयस चैट रूम और लाइव रूम मालिकों के लिए जिम्मेदारियाँ और अनुपालन आवश्यकताएँ।',
    vipDescPageTitle: 'VIP सदस्यता विवरण',
    vipDescPageDesc: 'Sara Live VIP स्तर, अपग्रेड नियम और प्रत्येक स्तर पर उपलब्ध विशेषाधिकारों के बारे में जानें।',
    lastUpdated: 'अंतिम अपडेट: 27 जुलाई 2026',
    cancelPageTitle: 'खाता हटाएँ',
    cancelPageDesc: 'अपना Sara Live खाता स्थायी रूप से हटाने के लिए निष्क्रियकरण कोड दर्ज करें। यह क्रिया वापस नहीं की जा सकती।',
    cancelCodeLabel: 'निष्क्रियकरण कोड',
    cancelCodePlaceholder: 'अपना निष्क्रियकरण कोड दर्ज करें',
    cancelCodeHint: 'अपना निष्क्रियकरण कोड Sara Live ऐप में खाता सेटिंग्स में देखें।',
    cancelSubmit: 'मेरा खाता हटाएँ',
    cancelConfirm: 'क्या आप वाकई इस खाते को स्थायी रूप से हटाना चाहते हैं? यह वापस नहीं किया जा सकता।',
    cancelSuccess: 'आपका खाता सफलतापूर्वक हटा दिया गया है।',
    cancelFail: 'खाता हटाया नहीं जा सका। कृपया कोड जाँचें और पुनः प्रयास करें।',
    cancelBlocked: 'बहुत अधिक असफल प्रयास। कृपया 2 घंटे बाद पुनः प्रयास करें।',
    cancelAlready: 'यह खाता पहले ही हटाया जा चुका है।',
    cancelNetworkError: 'नेटवर्क त्रुटि। कृपया बाद में पुनः प्रयास करें।',
    aboutPageTitle: 'Sara Live के बारे में',
    aboutP1: 'Sara Live एक लाइव-स्ट्रीमिंग और सोशल प्लेटफ़ॉर्म है जो केवल 18 वर्ष या उससे अधिक उम्र के उपयोगकर्ताओं के लिए है। यह दुनिया भर के वयस्कों को वास्तविक समय में जुड़ने, संवाद करने और अनुभव साझा करने में सक्षम बनाता है।',
    aboutP2: 'Sara Live लाइव वीडियो, वॉयस रूम, कॉल, सार्वजनिक चैट, डायरेक्ट मैसेज, शॉर्ट वीडियो, क्रिएटर टूल और वर्चुअल गिफ्ट प्रदान करता है। सुविधाओं की उपलब्धता देश, खाता स्थिति और ऐप संस्करण के अनुसार भिन्न हो सकती है।',
    aboutP3: 'हम एक सम्मानजनक और जिम्मेदार समुदाय बनाए रखने के लिए प्रतिबद्ध हैं। उपयोगकर्ता सामग्री या खातों की रिपोर्ट कर सकते हैं, अन्य उपयोगकर्ताओं को ब्लॉक कर सकते हैं, और सुरक्षा केंद्र में हमारे सामुदायिक दिशानिर्देश और बाल सुरक्षा मानक देख सकते हैं।',
    aboutSupportTitle: 'सहायता और संपर्क',
    aboutSupportDesc: 'Sara Live उस डेवलपर द्वारा संचालित है जो संबंधित ऐप-स्टोर लिस्टिंग में पहचाना जाता है। उत्पाद सहायता, सुरक्षा चिंताओं या कानूनी पूछताछ के लिए Sara Live में Help & Support का उपयोग करें। यदि इन-ऐप चैनल उपलब्ध नहीं है, तो उस ऐप-स्टोर लिस्टिंग में दिखाए गए डेवलपर संपर्क का उपयोग करें जहाँ से आपने Sara Live डाउनलोड किया था।',
    aboutThanks: 'Sara Live समुदाय का हिस्सा बनने के लिए धन्यवाद।',
    aboutSupportLink: 'सहायता',
  },
};

function getLang() {
  const params = new URLSearchParams(window.location.search);
  const fromQuery = params.get('lang');
  if (fromQuery && SITE_I18N[fromQuery]) return fromQuery;
  const saved = localStorage.getItem('sara-live-lang');
  if (saved && SITE_I18N[saved]) return saved;
  return 'en';
}

function setLang(lang) {
  if (!SITE_I18N[lang]) lang = 'en';
  localStorage.setItem('sara-live-lang', lang);
  document.documentElement.lang = lang === 'hi' ? 'hi' : lang;
  const t = SITE_I18N[lang];

  document.querySelectorAll('[data-i18n]').forEach((el) => {
    const key = el.getAttribute('data-i18n');
    if (t[key]) el.textContent = t[key];
  });

  document.querySelectorAll('[data-i18n-placeholder]').forEach((el) => {
    const key = el.getAttribute('data-i18n-placeholder');
    if (t[key]) el.setAttribute('placeholder', t[key]);
  });

  document.querySelectorAll('.lang-btn').forEach((btn) => {
    btn.classList.toggle('active', btn.dataset.lang === lang);
  });

  document.querySelectorAll('.lang-panel').forEach((panel) => {
    panel.classList.toggle('active', panel.dataset.lang === lang);
  });

  const privacyLinks = document.querySelectorAll('[data-privacy-email]');
  privacyLinks.forEach((link) => {
    link.href = 'mailto:' + APP_CONFIG.privacyEmail;
    link.textContent = APP_CONFIG.privacyEmail;
  });

  const supportLinks = document.querySelectorAll('[data-support-email]');
  supportLinks.forEach((link) => {
    link.href = 'mailto:' + APP_CONFIG.supportEmail;
    link.textContent = APP_CONFIG.supportEmail;
  });

  const titleEl = document.querySelector('[data-page-title-key]');
  if (titleEl) {
    const key = titleEl.getAttribute('data-page-title-key');
    document.title = (t[key] || APP_CONFIG.appName) + ' | ' + APP_CONFIG.appName;
  }

  if (document.body.dataset.legalPage === 'true') {
    const contentEl = document.getElementById('legal-content');
    const legalSection = document.body.dataset.legalSection || '';
    if (contentEl) {
      let html = '';
      if (legalSection === 'privacy' && typeof renderPrivacyContent === 'function') {
        html = renderPrivacyContent(lang);
      } else if (legalSection === 'tos' && typeof renderTermsContent === 'function') {
        html = renderTermsContent(lang);
      } else if (legalSection === 'creator-terms' && typeof renderCreatorTermsContent === 'function') {
        html = renderCreatorTermsContent(lang);
      } else if (legalSection === 'room-owner-terms' && typeof renderRoomOwnerTermsContent === 'function') {
        html = renderRoomOwnerTermsContent(lang);
      } else if (legalSection === 'vip-desc' && typeof renderVipDescContent === 'function') {
        html = renderVipDescContent(lang);
      } else if (typeof renderLegalDocument === 'function') {
        html = renderLegalDocument(lang, legalSection || undefined, { showTitle: false });
      }
      if (html) {
        contentEl.innerHTML = `<section class="legal-section" id="${legalSection || 'legal'}">${html}</section>`;
      }
    }
    if (!legalSection) {
      const tocEl = document.getElementById('legal-toc');
      if (tocEl && typeof renderLegalToc === 'function') {
        tocEl.innerHTML = renderLegalToc(lang);
        bindLegalTocLinks();
      }
      scrollToLegalHash();
    }
  }
}

function bindLegalTocLinks() {
  const panel = document.getElementById('toc-panel');
  const toggle = document.getElementById('toc-toggle');
  if (!panel) return;
  panel.querySelectorAll('a[href^="#"]').forEach((link) => {
    link.addEventListener('click', () => {
      if (!window.matchMedia('(max-width: 860px)').matches) return;
      panel.classList.remove('open');
      if (toggle) toggle.setAttribute('aria-expanded', 'false');
      closeMobileNav();
    });
  });
}

function closeMobileNav() {
  const navLinks = document.querySelector('.nav-links');
  const toggle = document.querySelector('.menu-toggle');
  if (navLinks) navLinks.classList.remove('open');
  if (toggle) toggle.setAttribute('aria-expanded', 'false');
  document.body.classList.remove('menu-open');
}

function initMobileNav() {
  const toggle = document.querySelector('.menu-toggle');
  const navLinks = document.querySelector('.nav-links');
  if (!toggle || !navLinks) return;

  toggle.setAttribute('aria-expanded', 'false');
  toggle.addEventListener('click', () => {
    const open = navLinks.classList.toggle('open');
    toggle.setAttribute('aria-expanded', open ? 'true' : 'false');
    document.body.classList.toggle('menu-open', open);
  });

  navLinks.querySelectorAll('a').forEach((link) => {
    link.addEventListener('click', () => closeMobileNav());
  });

  document.addEventListener('keydown', (event) => {
    if (event.key === 'Escape') closeMobileNav();
  });
}

function initLegalToc() {
  const toggle = document.getElementById('toc-toggle');
  const panel = document.getElementById('toc-panel');
  if (!toggle || !panel) return;

  toggle.addEventListener('click', () => {
    const open = panel.classList.toggle('open');
    toggle.setAttribute('aria-expanded', open ? 'true' : 'false');
  });

  if (window.matchMedia('(min-width: 861px)').matches) {
    panel.classList.add('open');
    toggle.setAttribute('aria-expanded', 'true');
  }
}

function getApiBaseUrl() {
  const base = (APP_CONFIG.apiBaseUrl || '').trim();
  if (base) {
    return base.replace(/\/+$/, '');
  }
  return window.location.origin;
}

function showCancelMessage(text, type) {
  const messageEl = document.getElementById('cancel-message');
  if (!messageEl) return;
  messageEl.hidden = false;
  messageEl.textContent = text;
  messageEl.className = 'cancel-message cancel-message--' + (type || 'info');
}

async function submitCancelAccount(cancelCode) {
  const lang = getLang();
  const t = SITE_I18N[lang] || SITE_I18N.en;
  const url = getApiBaseUrl() + '/userInfo/cancelAccountByCode';
  const response = await fetch(url, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ cancelCode }),
  });
  const payload = await response.json().catch(() => ({}));
  if (!response.ok) {
    throw new Error(t.cancelNetworkError);
  }
  if (payload.code === 0 && payload.data && payload.data.success) {
    return t.cancelSuccess;
  }
  if (payload.code === 8) {
    throw new Error(t.cancelAlready);
  }
  if (payload.code === 91) {
    throw new Error(t.cancelBlocked);
  }
  throw new Error(t.cancelFail);
}

function initCancelAccountForm() {
  const form = document.getElementById('cancel-account-form');
  if (!form || form.dataset.bound === 'true') return;
  form.dataset.bound = 'true';

  form.addEventListener('submit', async (event) => {
    event.preventDefault();
    const lang = getLang();
    const t = SITE_I18N[lang] || SITE_I18N.en;
    const input = document.getElementById('cancel-code');
    const submitBtn = form.querySelector('button[type="submit"]');
    const cancelCode = (input && input.value ? input.value : '').trim();
    if (!cancelCode) {
      showCancelMessage(t.cancelFail, 'error');
      return;
    }
    if (!window.confirm(t.cancelConfirm)) {
      return;
    }

    if (submitBtn) submitBtn.disabled = true;
    showCancelMessage('', 'info');
    try {
      const message = await submitCancelAccount(cancelCode);
      showCancelMessage(message, 'success');
      form.reset();
    } catch (error) {
      showCancelMessage(error.message || t.cancelFail, 'error');
    } finally {
      if (submitBtn) submitBtn.disabled = false;
    }
  });
}

function scrollToLegalHash() {
  if (document.body.dataset.legalSection) {
    return;
  }
  const hash = window.location.hash || '#privacy';
  const target = document.querySelector(hash);
  if (target) {
    target.scrollIntoView({ behavior: 'smooth', block: 'start' });
  }
}

function redirectLegacyLegalUrls() {
  const page = (window.location.pathname.split('/').pop() || '').toLowerCase();
  const hash = window.location.hash;
  if (page === 'privacy.html' && hash === '#tos') {
    window.location.replace('terms.html' + window.location.search);
  }
}

function initAboutPage() {
  if (document.body.dataset.aboutPage !== 'true') return;
  initStandaloneBackButton();
}

function initSafetyCenterPage() {
  if (document.body.dataset.safetyCenterPage !== 'true') return;
  initStandaloneBackButton();
}

function initStandaloneBackButton() {
  const appNavigation = window.SoraLegalNavigation;
  if (appNavigation?.postMessage) {
    document.documentElement.classList.add('app-embedded');
  }

  const goBack = () => {
    if (document.referrer && window.history.length > 1) {
      window.history.back();
      return;
    }
    if (appNavigation?.postMessage) {
      appNavigation.postMessage('back');
      return;
    }
    const fallback = document.body.dataset.backFallback;
    if (fallback) {
      window.location.assign(fallback);
    }
  };

  document.querySelectorAll('[data-back-button]').forEach((button) => {
    button.addEventListener('click', goBack);
  });
}

function initSite() {
  redirectLegacyLegalUrls();
  const lang = getLang();
  setLang(lang);

  document.querySelectorAll('.lang-btn').forEach((btn) => {
    btn.addEventListener('click', () => setLang(btn.dataset.lang));
  });

  initMobileNav();
  initLegalToc();
  initCancelAccountForm();
  initAboutPage();
  initSafetyCenterPage();

  const current = (window.location.pathname.split('/').pop() || 'index.html').toLowerCase();
  document.querySelectorAll('.nav-links a[data-nav]').forEach((link) => {
    const href = link.getAttribute('href').toLowerCase();
    const isHome = current === 'index.html' || current === '';
    link.classList.toggle('active', (isHome && link.dataset.nav === 'home') || href.includes(current.replace('.html', '')));
  });

  window.addEventListener('hashchange', scrollToLegalHash);
}

document.addEventListener('DOMContentLoaded', initSite);
