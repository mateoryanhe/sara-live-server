// ignore: unused_import
import 'package:intl/intl.dart' as intl;
import 'app_localizations.dart';

// ignore_for_file: type=lint

/// The translations for English (`en`).
class AppLocalizationsEn extends AppLocalizations {
  AppLocalizationsEn([String locale = 'en']) : super(locale);

  @override
  String get appTitle => 'Sara Live';

  @override
  String get uiRenderErrorTitle => 'UI render error';

  @override
  String get commonCancel => 'Cancel';

  @override
  String get commonSave => 'Save';

  @override
  String get commonRetry => 'Retry';

  @override
  String get commonCopied => 'Copied';

  @override
  String get inviteCodeCopied => 'Invite code copied';

  @override
  String get languageTitle => 'Language';

  @override
  String get languageSubtitle => 'English or Chinese';

  @override
  String get languageEnglish => 'English';

  @override
  String get languageChinese => '中文';

  @override
  String get languagePickerTitle => 'Language';

  @override
  String get languagePickerHint => 'Your choice is saved automatically';

  @override
  String get uiStyleTitle => 'Appearance';

  @override
  String get uiStyleSubtitle => 'Tap to change theme';

  @override
  String get settingsSection => 'Settings';

  @override
  String get accountSection => 'Account';

  @override
  String get logout => 'Log out';

  @override
  String get editNicknameTitle => 'Edit nickname';

  @override
  String get editNicknameHint => 'Enter a new nickname';

  @override
  String get avatarFormatError => 'Only jpg, png, webp, gif, or bmp';

  @override
  String userDisplayName(int userId) {
    return 'User $userId';
  }

  @override
  String get walletGold => 'Coins';

  @override
  String get walletDiamond => 'Diamonds';

  @override
  String get walletFollow => 'Following';

  @override
  String get walletFans => 'Followers';

  @override
  String get vipTitle => 'VIP';

  @override
  String get inviteCodeTitle => 'Invite code';

  @override
  String get anchorRoleTitle => 'Host status';

  @override
  String get anchorRoleEnabled => 'Enabled';

  @override
  String get anchorRoleDisabled => 'Not enabled';

  @override
  String get tabOneToOne => '1v1';

  @override
  String get tabLive => 'Live';

  @override
  String get tabGame => 'Games';

  @override
  String get tabMine => 'Me';

  @override
  String get tabOnlineAudience => '1v1';

  @override
  String get tabOneToOneRoom => '1v1 room';

  @override
  String get tabShowLive => 'Go live';

  @override
  String get placeholderLiveTitle => 'Live';

  @override
  String get placeholderLiveSubtitle =>
      'Popular show and game streams coming soon';

  @override
  String get placeholderGameTitle => 'Games';

  @override
  String get placeholderGameSubtitle =>
      'Game catalog and mini-games coming soon';

  @override
  String get deviceCodeTitle => 'Device ID';

  @override
  String get deviceCodeEditTitle => 'Edit device ID';

  @override
  String get deviceCodeHint => 'Used for quick device login (max 128 chars)';

  @override
  String get deviceCodeWarning =>
      'Saving will log in with this ID. You may switch to another account.';

  @override
  String get deviceCodeEmpty => 'Device ID cannot be empty';

  @override
  String get deviceCodeTooLong => 'Device ID cannot exceed 128 characters';

  @override
  String get deviceCodeUnchanged => 'Device ID unchanged';

  @override
  String get deviceCodeSwitchSuccess => 'Device ID updated; logged in';

  @override
  String get deviceCodeNotSet => 'Not set (auto-generated on login)';

  @override
  String get apiBaseTitle => 'API server';

  @override
  String get apiBaseEditTitle => 'Edit API base URL';

  @override
  String get apiBaseHint => 'e.g. http://192.168.2.11:9443';

  @override
  String get apiBaseWarning =>
      'Saving switches the server and logs you out. You will sign in again on the new host.';

  @override
  String get apiBaseResetDefault => 'Use default';

  @override
  String get apiBaseEmpty => 'API URL cannot be empty';

  @override
  String get apiBaseInvalid => 'Enter a valid http or https URL';

  @override
  String get apiBaseUnchanged => 'API server unchanged';

  @override
  String get apiBaseSwitchSuccess => 'API server updated; please sign in again';

  @override
  String get apiBaseNotSet => 'Not configured';

  @override
  String apiBaseUsingDefault(String url) {
    return 'Default: $url';
  }

  @override
  String get liveHotListTitle => 'Live now';

  @override
  String get liveHotListEmpty => 'No live rooms right now';

  @override
  String get liveBadgeOnAir => 'LIVE';

  @override
  String liveOnlineCount(int count) {
    return '$count watching';
  }

  @override
  String get liveChatHint => 'Say something…';

  @override
  String get liveSendGift => 'Send gift';

  @override
  String get liveGiftListEmpty => 'No gifts available';

  @override
  String get liveAnchorEnded => 'The host ended the live stream';

  @override
  String get liveWatchBlockedInCall => 'Exit the call before watching live';

  @override
  String get liveJoiningRoom => 'Joining live room…';

  @override
  String get liveWatchAudienceTitle => 'Viewers online';

  @override
  String liveWatchAudienceSubtitle(int count) {
    return '$count watching (excluding host)';
  }

  @override
  String get liveWatchAudienceEmpty => 'No viewers yet';

  @override
  String get liveWatchTapToChat => 'Chat';

  @override
  String liveGiftSent(String giftName, double diamond) {
    return 'Sent $giftName. Balance: $diamond';
  }
}
