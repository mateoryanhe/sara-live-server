import 'package:cached_network_image/cached_network_image.dart';
import 'package:flutter/material.dart';

import '../../core/cache/image_disk_cache.dart';
import '../../data/live_room/live_room_models.dart';
import '../../data/vip/vip_chat_appearance.dart';
import '../../data/vip/vip_chat_bubble_skin.dart';
import '../../data/vip/vip_privilege_frames.dart';
import '../../l10n/app_localizations.dart';
import 'live_room_chat_avatar.dart';

/// 公屏一行：三列（头像 | VIP 等级标识 | 文案）；VIP0 与其它等级同一结构。
class LiveRoomChatLineTile extends StatelessWidget {
  const LiveRoomChatLineTile({
    super.key,
    required this.line,
    required this.vipFrames,
  });

  final LiveRoomChatLine line;
  final VipPrivilegeFrames vipFrames;

  static const _avatarSize = 30.0;
  static const _badgeWidth = 58.0;

  @override
  Widget build(BuildContext context) {
    final name = line.senderName.isNotEmpty ? line.senderName : '…';
    final vipLevel = line.vipLevel;
    final appearance = VipChatAppearance.forLevel(vipLevel);

    if (line.isGift) {
      if (line.giftName.isEmpty && line.giftIcon.isEmpty) {
        return const SizedBox.shrink();
      }
    } else if (line.content.isEmpty) {
      return const SizedBox.shrink();
    }

    return SizedBox(
      width: double.infinity,
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.center,
        children: [
          LiveRoomChatAvatar(
          avatarUrl: line.senderAvatar,
          displayName: name,
          frameUrl: vipFrames.chatAvatarFrameUrl(vipLevel),
          size: _avatarSize,
        ),
        const SizedBox(width: 6),
        SizedBox(
          width: _badgeWidth,
          child: Align(
            alignment: Alignment.centerLeft,
            child: _MiddleBadge(
              vipLevel: vipLevel,
              levelIconUrl: vipFrames.levelIconUrl(vipLevel),
              levelLabel: vipFrames.levelLabel(vipLevel),
              appearance: appearance,
            ),
          ),
        ),
        const SizedBox(width: 4),
        Expanded(
          child: Align(
            alignment: Alignment.centerLeft,
            child: line.isGift
                ? _GiftMessageBubble(
                    giftName: line.giftName,
                    giftIcon: line.giftIcon,
                    giftCount: line.giftCount,
                    actionLabel:
                        AppLocalizations.of(context).liveChatGiftAction,
                    appearance: appearance,
                  )
                : _MessageBubble(
                    text: line.content,
                    appearance: appearance,
                    bubbleFrameUrl:
                        vipFrames.chatMessageBubbleFrameUrl(vipLevel),
                  ),
          ),
        ),
      ],
      ),
    );
  }
}

class _MiddleBadge extends StatelessWidget {
  const _MiddleBadge({
    required this.vipLevel,
    required this.levelIconUrl,
    required this.levelLabel,
    required this.appearance,
  });

  final int vipLevel;
  final String? levelIconUrl;
  final String levelLabel;
  final VipChatAppearance appearance;

  @override
  Widget build(BuildContext context) {
    final icon = levelIconUrl;
    if (icon != null) {
      return CachedNetworkImage(
        imageUrl: icon,
        cacheManager: ImageDiskCache.instance.manager,
        height: 22,
        fit: BoxFit.contain,
        alignment: Alignment.centerLeft,
        errorWidget: (_, __, ___) => _GradientVipPill(
          vipLevel: vipLevel,
          label: levelLabel,
          appearance: appearance,
        ),
      );
    }
    return _GradientVipPill(
      vipLevel: vipLevel,
      label: levelLabel,
      appearance: appearance,
    );
  }
}

class _GradientVipPill extends StatelessWidget {
  const _GradientVipPill({
    required this.vipLevel,
    required this.label,
    required this.appearance,
  });

  final int vipLevel;
  final String label;
  final VipChatAppearance appearance;

  @override
  Widget build(BuildContext context) {
    final isVip0 = vipLevel <= 0;
    return DecoratedBox(
      decoration: BoxDecoration(
        borderRadius: BorderRadius.circular(999),
        gradient: appearance.badgeGradient,
        boxShadow: isVip0
            ? null
            : [
                BoxShadow(
                  color: Colors.black.withValues(alpha: 0.35),
                  blurRadius: 4,
                  offset: const Offset(0, 1),
                ),
              ],
      ),
      child: Padding(
        padding: const EdgeInsets.fromLTRB(6, 3, 8, 3),
        child: Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(
              isVip0 ? Icons.person_rounded : Icons.workspace_premium_rounded,
              size: 13,
              color: appearance.badgeTextColor,
            ),
            const SizedBox(width: 3),
            Flexible(
              child: Text(
                label,
                maxLines: 1,
                overflow: TextOverflow.ellipsis,
                style: TextStyle(
                  color: appearance.badgeTextColor,
                  fontSize: 10,
                  fontWeight: FontWeight.w800,
                  height: 1.1,
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _GiftMessageBubble extends StatelessWidget {
  const _GiftMessageBubble({
    required this.giftName,
    required this.giftIcon,
    required this.giftCount,
    required this.actionLabel,
    required this.appearance,
  });

  final String giftName;
  final String giftIcon;
  final int giftCount;
  final String actionLabel;
  final VipChatAppearance appearance;

  @override
  Widget build(BuildContext context) {
    final name = giftName.isNotEmpty ? giftName : 'Gift';
    final count = giftCount > 0 ? giftCount : 1;

    return DecoratedBox(
      decoration: BoxDecoration(
        borderRadius: BorderRadius.circular(12),
        gradient: const LinearGradient(
          colors: [Color(0xFFFBBF24), Color(0xFFE11D48)],
        ),
        boxShadow: [
          BoxShadow(
            color: const Color(0xFFF59E0B).withValues(alpha: 0.35),
            blurRadius: 8,
          ),
        ],
      ),
      child: Padding(
        padding: const EdgeInsets.all(1.4),
        child: Container(
          padding: const EdgeInsets.fromLTRB(8, 6, 8, 6),
          decoration: BoxDecoration(
            color: const Color(0xCC2A1060),
            borderRadius: BorderRadius.circular(10.5),
          ),
          child: Row(
            children: [
              _GiftIconThumb(url: giftIcon),
              const SizedBox(width: 8),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Text(
                      actionLabel,
                      style: TextStyle(
                        color: Colors.white.withValues(alpha: 0.65),
                        fontSize: 10,
                        fontWeight: FontWeight.w600,
                        height: 1.1,
                      ),
                    ),
                    Text(
                      name,
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                      style: const TextStyle(
                        color: Color(0xFFFFE566),
                        fontSize: 13,
                        fontWeight: FontWeight.w800,
                        height: 1.25,
                      ),
                    ),
                  ],
                ),
              ),
              const SizedBox(width: 6),
              _GiftCountBadge(count: count),
            ],
          ),
        ),
      ),
    );
  }
}

class _GiftIconThumb extends StatelessWidget {
  const _GiftIconThumb({required this.url});

  final String url;

  @override
  Widget build(BuildContext context) {
    return Container(
      width: 40,
      height: 40,
      decoration: BoxDecoration(
        color: Colors.white.withValues(alpha: 0.12),
        borderRadius: BorderRadius.circular(10),
        border: Border.all(color: Colors.white24),
      ),
      clipBehavior: Clip.antiAlias,
      child: url.trim().isEmpty
          ? const Icon(Icons.card_giftcard_rounded, color: Color(0xFFFFD166))
          : CachedNetworkImage(
              imageUrl: url.trim(),
              cacheManager: ImageDiskCache.instance.manager,
              fit: BoxFit.cover,
              errorWidget: (_, __, ___) => const Icon(
                Icons.card_giftcard_rounded,
                color: Color(0xFFFFD166),
              ),
            ),
    );
  }
}

class _GiftCountBadge extends StatelessWidget {
  const _GiftCountBadge({required this.count});

  final int count;

  @override
  Widget build(BuildContext context) {
    return Container(
      constraints: const BoxConstraints(minWidth: 28),
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
      decoration: BoxDecoration(
        gradient: const LinearGradient(
          colors: [Color(0xFFFFE566), Color(0xFFF59E0B)],
        ),
        borderRadius: BorderRadius.circular(8),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withValues(alpha: 0.25),
            blurRadius: 4,
            offset: const Offset(0, 1),
          ),
        ],
      ),
      child: Text(
        '×$count',
        textAlign: TextAlign.center,
        style: const TextStyle(
          color: Color(0xFF422006),
          fontSize: 13,
          fontWeight: FontWeight.w900,
          height: 1,
        ),
      ),
    );
  }
}

class _MessageBubble extends StatelessWidget {
  const _MessageBubble({
    required this.text,
    required this.appearance,
    this.bubbleFrameUrl,
  });

  final String text;
  final VipChatAppearance appearance;
  final String? bubbleFrameUrl;

  @override
  Widget build(BuildContext context) {
    final frame = bubbleFrameUrl?.trim() ?? '';
    if (frame.isNotEmpty) {
      return _NinePatchBubbleMessage(text: text, frameSource: frame);
    }

    final textStyle = TextStyle(
      color: appearance.textColor,
      fontSize: 13,
      height: 1.35,
      fontWeight: FontWeight.w500,
    );

    return IntrinsicWidth(
      child: DecoratedBox(
        decoration: BoxDecoration(
          borderRadius: BorderRadius.circular(12),
          gradient: appearance.borderGradient,
          boxShadow: [
            BoxShadow(
              color: appearance.borderGradient.colors.first
                  .withValues(alpha: 0.35),
              blurRadius: 6,
              spreadRadius: 0,
            ),
          ],
        ),
        child: Padding(
          padding: const EdgeInsets.all(1.4),
          child: Container(
            padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 6),
            decoration: BoxDecoration(
              color: appearance.fillColor,
              borderRadius: BorderRadius.circular(10.5),
            ),
            child: Text(text, style: textStyle),
          ),
        ),
      ),
    );
  }
}

/// VIP 公屏文字：点九气泡框（默认 APK 内置图，CMS 可配网络图）+ 文案。
class _NinePatchBubbleMessage extends StatelessWidget {
  const _NinePatchBubbleMessage({
    required this.text,
    required this.frameSource,
  });

  final String text;
  final String frameSource;

  @override
  Widget build(BuildContext context) {
    final provider = VipChatBubbleSkin.frameImageProvider(frameSource);
    const textStyle = TextStyle(
      color: Colors.white,
      fontSize: 13,
      height: 1.35,
      fontWeight: FontWeight.w600,
      shadows: [
        Shadow(
          color: Color(0x99000000),
          blurRadius: 3,
          offset: Offset(0, 1),
        ),
      ],
    );

    // 用 Container + DecorationImage 点九（比 Stack+Positioned.fill 在 Android 上更稳）。
    return IntrinsicWidth(
      child: _NinePatchBubbleFrame(
        image: provider,
        padding: VipChatBubbleSkin.contentPadding,
        child: Text(text, style: textStyle),
      ),
    );
  }
}

/// 点九气泡底图；加载失败时仍保留圆角底色，避免「只有裸字没有框」。
class _NinePatchBubbleFrame extends StatelessWidget {
  const _NinePatchBubbleFrame({
    required this.image,
    required this.padding,
    required this.child,
  });

  final ImageProvider image;
  final EdgeInsets padding;
  final Widget child;

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: padding,
      decoration: BoxDecoration(
        color: Colors.black.withValues(alpha: 0.42),
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: Colors.white24),
        image: DecorationImage(
          image: image,
          fit: BoxFit.fill,
          centerSlice: VipChatBubbleSkin.centerSlice,
          filterQuality: FilterQuality.medium,
        ),
      ),
      child: child,
    );
  }
}
