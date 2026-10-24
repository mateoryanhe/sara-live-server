import 'dart:async';

import 'package:flutter/foundation.dart';
import 'package:flutter/widgets.dart';
import 'package:flutter_localizations/flutter_localizations.dart';
import 'package:intl/intl.dart' as intl;

import 'app_localizations_en.dart';
import 'app_localizations_zh.dart';

// ignore_for_file: type=lint

/// Callers can lookup localized strings with an instance of AppLocalizations
/// returned by `AppLocalizations.of(context)`.
///
/// Applications need to include `AppLocalizations.delegate()` in their app's
/// `localizationDelegates` list, and the locales they support in the app's
/// `supportedLocales` list. For example:
///
/// ```dart
/// import 'l10n/app_localizations.dart';
///
/// return MaterialApp(
///   localizationsDelegates: AppLocalizations.localizationsDelegates,
///   supportedLocales: AppLocalizations.supportedLocales,
///   home: MyApplicationHome(),
/// );
/// ```
///
/// ## Update pubspec.yaml
///
/// Please make sure to update your pubspec.yaml to include the following
/// packages:
///
/// ```yaml
/// dependencies:
///   # Internationalization support.
///   flutter_localizations:
///     sdk: flutter
///   intl: any # Use the pinned version from flutter_localizations
///
///   # Rest of dependencies
/// ```
///
/// ## iOS Applications
///
/// iOS applications define key application metadata, including supported
/// locales, in an Info.plist file that is built into the application bundle.
/// To configure the locales supported by your app, you’ll need to edit this
/// file.
///
/// First, open your project’s ios/Runner.xcworkspace Xcode workspace file.
/// Then, in the Project Navigator, open the Info.plist file under the Runner
/// project’s Runner folder.
///
/// Next, select the Information Property List item, select Add Item from the
/// Editor menu, then select Localizations from the pop-up menu.
///
/// Select and expand the newly-created Localizations item then, for each
/// locale your application supports, add a new item and select the locale
/// you wish to add from the pop-up menu in the Value field. This list should
/// be consistent with the languages listed in the AppLocalizations.supportedLocales
/// property.
abstract class AppLocalizations {
  AppLocalizations(String locale)
      : localeName = intl.Intl.canonicalizedLocale(locale.toString());

  final String localeName;

  static AppLocalizations of(BuildContext context) {
    return Localizations.of<AppLocalizations>(context, AppLocalizations)!;
  }

  static const LocalizationsDelegate<AppLocalizations> delegate =
      _AppLocalizationsDelegate();

  /// A list of this localizations delegate along with the default localizations
  /// delegates.
  ///
  /// Returns a list of localizations delegates containing this delegate along with
  /// GlobalMaterialLocalizations.delegate, GlobalCupertinoLocalizations.delegate,
  /// and GlobalWidgetsLocalizations.delegate.
  ///
  /// Additional delegates can be added by appending to this list in
  /// MaterialApp. This list does not have to be used at all if a custom list
  /// of delegates is preferred or required.
  static const List<LocalizationsDelegate<dynamic>> localizationsDelegates =
      <LocalizationsDelegate<dynamic>>[
    delegate,
    GlobalMaterialLocalizations.delegate,
    GlobalCupertinoLocalizations.delegate,
    GlobalWidgetsLocalizations.delegate,
  ];

  /// A list of this localizations delegate's supported locales.
  static const List<Locale> supportedLocales = <Locale>[
    Locale('en'),
    Locale('zh')
  ];

  /// No description provided for @appTitle.
  ///
  /// In en, this message translates to:
  /// **'Sara Live'**
  String get appTitle;

  /// No description provided for @uiRenderErrorTitle.
  ///
  /// In en, this message translates to:
  /// **'UI render error'**
  String get uiRenderErrorTitle;

  /// No description provided for @commonCancel.
  ///
  /// In en, this message translates to:
  /// **'Cancel'**
  String get commonCancel;

  /// No description provided for @commonSave.
  ///
  /// In en, this message translates to:
  /// **'Save'**
  String get commonSave;

  /// No description provided for @commonRetry.
  ///
  /// In en, this message translates to:
  /// **'Retry'**
  String get commonRetry;

  /// No description provided for @commonCopied.
  ///
  /// In en, this message translates to:
  /// **'Copied'**
  String get commonCopied;

  /// No description provided for @inviteCodeCopied.
  ///
  /// In en, this message translates to:
  /// **'Invite code copied'**
  String get inviteCodeCopied;

  /// No description provided for @languageTitle.
  ///
  /// In en, this message translates to:
  /// **'Language'**
  String get languageTitle;

  /// No description provided for @languageSubtitle.
  ///
  /// In en, this message translates to:
  /// **'English or Chinese'**
  String get languageSubtitle;

  /// No description provided for @languageEnglish.
  ///
  /// In en, this message translates to:
  /// **'English'**
  String get languageEnglish;

  /// No description provided for @languageChinese.
  ///
  /// In en, this message translates to:
  /// **'中文'**
  String get languageChinese;

  /// No description provided for @languagePickerTitle.
  ///
  /// In en, this message translates to:
  /// **'Language'**
  String get languagePickerTitle;

  /// No description provided for @languagePickerHint.
  ///
  /// In en, this message translates to:
  /// **'Your choice is saved automatically'**
  String get languagePickerHint;

  /// No description provided for @uiStyleTitle.
  ///
  /// In en, this message translates to:
  /// **'Appearance'**
  String get uiStyleTitle;

  /// No description provided for @uiStyleSubtitle.
  ///
  /// In en, this message translates to:
  /// **'Tap to change theme'**
  String get uiStyleSubtitle;

  /// No description provided for @settingsSection.
  ///
  /// In en, this message translates to:
  /// **'Settings'**
  String get settingsSection;

  /// No description provided for @accountSection.
  ///
  /// In en, this message translates to:
  /// **'Account'**
  String get accountSection;

  /// No description provided for @logout.
  ///
  /// In en, this message translates to:
  /// **'Log out'**
  String get logout;

  /// No description provided for @editNicknameTitle.
  ///
  /// In en, this message translates to:
  /// **'Edit nickname'**
  String get editNicknameTitle;

  /// No description provided for @editNicknameHint.
  ///
  /// In en, this message translates to:
  /// **'Enter a new nickname'**
  String get editNicknameHint;

  /// No description provided for @avatarFormatError.
  ///
  /// In en, this message translates to:
  /// **'Only jpg, png, webp, gif, or bmp'**
  String get avatarFormatError;

  /// No description provided for @userDisplayName.
  ///
  /// In en, this message translates to:
  /// **'User {userId}'**
  String userDisplayName(int userId);

  /// No description provided for @walletGold.
  ///
  /// In en, this message translates to:
  /// **'Coins'**
  String get walletGold;

  /// No description provided for @walletDiamond.
  ///
  /// In en, this message translates to:
  /// **'Diamonds'**
  String get walletDiamond;

  /// No description provided for @walletFollow.
  ///
  /// In en, this message translates to:
  /// **'Following'**
  String get walletFollow;

  /// No description provided for @walletFans.
  ///
  /// In en, this message translates to:
  /// **'Followers'**
  String get walletFans;

  /// No description provided for @vipTitle.
  ///
  /// In en, this message translates to:
  /// **'VIP'**
  String get vipTitle;

  /// No description provided for @inviteCodeTitle.
  ///
  /// In en, this message translates to:
  /// **'Invite code'**
  String get inviteCodeTitle;

  /// No description provided for @anchorRoleTitle.
  ///
  /// In en, this message translates to:
  /// **'Host status'**
  String get anchorRoleTitle;

  /// No description provided for @anchorRoleEnabled.
  ///
  /// In en, this message translates to:
  /// **'Enabled'**
  String get anchorRoleEnabled;

  /// No description provided for @anchorRoleDisabled.
  ///
  /// In en, this message translates to:
  /// **'Not enabled'**
  String get anchorRoleDisabled;

  /// No description provided for @tabOneToOne.
  ///
  /// In en, this message translates to:
  /// **'1v1'**
  String get tabOneToOne;

  /// No description provided for @tabLive.
  ///
  /// In en, this message translates to:
  /// **'Live'**
  String get tabLive;

  /// No description provided for @tabGame.
  ///
  /// In en, this message translates to:
  /// **'Games'**
  String get tabGame;

  /// No description provided for @tabMine.
  ///
  /// In en, this message translates to:
  /// **'Me'**
  String get tabMine;

  /// No description provided for @tabOnlineAudience.
  ///
  /// In en, this message translates to:
  /// **'1v1'**
  String get tabOnlineAudience;

  /// No description provided for @tabOneToOneRoom.
  ///
  /// In en, this message translates to:
  /// **'1v1 room'**
  String get tabOneToOneRoom;

  /// No description provided for @tabShowLive.
  ///
  /// In en, this message translates to:
  /// **'Go live'**
  String get tabShowLive;

  /// No description provided for @placeholderLiveTitle.
  ///
  /// In en, this message translates to:
  /// **'Live'**
  String get placeholderLiveTitle;

  /// No description provided for @placeholderLiveSubtitle.
  ///
  /// In en, this message translates to:
  /// **'Popular show and game streams coming soon'**
  String get placeholderLiveSubtitle;

  /// No description provided for @placeholderGameTitle.
  ///
  /// In en, this message translates to:
  /// **'Games'**
  String get placeholderGameTitle;

  /// No description provided for @placeholderGameSubtitle.
  ///
  /// In en, this message translates to:
  /// **'Game catalog and mini-games coming soon'**
  String get placeholderGameSubtitle;

  /// No description provided for @deviceCodeTitle.
  ///
  /// In en, this message translates to:
  /// **'Device ID'**
  String get deviceCodeTitle;

  /// No description provided for @deviceCodeEditTitle.
  ///
  /// In en, this message translates to:
  /// **'Edit device ID'**
  String get deviceCodeEditTitle;

  /// No description provided for @deviceCodeHint.
  ///
  /// In en, this message translates to:
  /// **'Used for quick device login (max 128 chars)'**
  String get deviceCodeHint;

  /// No description provided for @deviceCodeWarning.
  ///
  /// In en, this message translates to:
  /// **'Saving will log in with this ID. You may switch to another account.'**
  String get deviceCodeWarning;

  /// No description provided for @deviceCodeEmpty.
  ///
  /// In en, this message translates to:
  /// **'Device ID cannot be empty'**
  String get deviceCodeEmpty;

  /// No description provided for @deviceCodeTooLong.
  ///
  /// In en, this message translates to:
  /// **'Device ID cannot exceed 128 characters'**
  String get deviceCodeTooLong;

  /// No description provided for @deviceCodeUnchanged.
  ///
  /// In en, this message translates to:
  /// **'Device ID unchanged'**
  String get deviceCodeUnchanged;

  /// No description provided for @deviceCodeSwitchSuccess.
  ///
  /// In en, this message translates to:
  /// **'Device ID updated; logged in'**
  String get deviceCodeSwitchSuccess;

  /// No description provided for @deviceCodeNotSet.
  ///
  /// In en, this message translates to:
  /// **'Not set (auto-generated on login)'**
  String get deviceCodeNotSet;

  /// No description provided for @apiBaseTitle.
  ///
  /// In en, this message translates to:
  /// **'API server'**
  String get apiBaseTitle;

  /// No description provided for @apiBaseEditTitle.
  ///
  /// In en, this message translates to:
  /// **'Edit API base URL'**
  String get apiBaseEditTitle;

  /// No description provided for @apiBaseHint.
  ///
  /// In en, this message translates to:
  /// **'e.g. http://192.168.2.11:9443'**
  String get apiBaseHint;

  /// No description provided for @apiBaseWarning.
  ///
  /// In en, this message translates to:
  /// **'Saving switches the server and logs you out. You will sign in again on the new host.'**
  String get apiBaseWarning;

  /// No description provided for @apiBaseResetDefault.
  ///
  /// In en, this message translates to:
  /// **'Use default'**
  String get apiBaseResetDefault;

  /// No description provided for @apiBaseEmpty.
  ///
  /// In en, this message translates to:
  /// **'API URL cannot be empty'**
  String get apiBaseEmpty;

  /// No description provided for @apiBaseInvalid.
  ///
  /// In en, this message translates to:
  /// **'Enter a valid http or https URL'**
  String get apiBaseInvalid;

  /// No description provided for @apiBaseUnchanged.
  ///
  /// In en, this message translates to:
  /// **'API server unchanged'**
  String get apiBaseUnchanged;

  /// No description provided for @apiBaseSwitchSuccess.
  ///
  /// In en, this message translates to:
  /// **'API server updated; please sign in again'**
  String get apiBaseSwitchSuccess;

  /// No description provided for @apiBaseNotSet.
  ///
  /// In en, this message translates to:
  /// **'Not configured'**
  String get apiBaseNotSet;

  /// No description provided for @apiBaseUsingDefault.
  ///
  /// In en, this message translates to:
  /// **'Default: {url}'**
  String apiBaseUsingDefault(String url);

  /// No description provided for @liveHotListTitle.
  ///
  /// In en, this message translates to:
  /// **'Live now'**
  String get liveHotListTitle;

  /// No description provided for @liveHotListEmpty.
  ///
  /// In en, this message translates to:
  /// **'No live rooms right now'**
  String get liveHotListEmpty;

  /// No description provided for @liveBadgeOnAir.
  ///
  /// In en, this message translates to:
  /// **'LIVE'**
  String get liveBadgeOnAir;

  /// No description provided for @liveOnlineCount.
  ///
  /// In en, this message translates to:
  /// **'{count} watching'**
  String liveOnlineCount(int count);

  /// No description provided for @liveChatHint.
  ///
  /// In en, this message translates to:
  /// **'Say something…'**
  String get liveChatHint;

  /// No description provided for @liveSendGift.
  ///
  /// In en, this message translates to:
  /// **'Send gift'**
  String get liveSendGift;

  /// No description provided for @liveGiftListEmpty.
  ///
  /// In en, this message translates to:
  /// **'No gifts available'**
  String get liveGiftListEmpty;

  /// No description provided for @liveAnchorEnded.
  ///
  /// In en, this message translates to:
  /// **'The host ended the live stream'**
  String get liveAnchorEnded;

  /// No description provided for @liveWatchBlockedInCall.
  ///
  /// In en, this message translates to:
  /// **'Exit the call before watching live'**
  String get liveWatchBlockedInCall;

  /// No description provided for @liveJoiningRoom.
  ///
  /// In en, this message translates to:
  /// **'Joining live room…'**
  String get liveJoiningRoom;

  /// No description provided for @liveWatchAudienceTitle.
  ///
  /// In en, this message translates to:
  /// **'Viewers online'**
  String get liveWatchAudienceTitle;

  /// No description provided for @liveWatchAudienceSubtitle.
  ///
  /// In en, this message translates to:
  /// **'{count} watching (excluding host)'**
  String liveWatchAudienceSubtitle(int count);

  /// No description provided for @liveWatchAudienceEmpty.
  ///
  /// In en, this message translates to:
  /// **'No viewers yet'**
  String get liveWatchAudienceEmpty;

  /// No description provided for @liveWatchTapToChat.
  ///
  /// In en, this message translates to:
  /// **'Chat'**
  String get liveWatchTapToChat;

  /// No description provided for @liveGiftSent.
  ///
  /// In en, this message translates to:
  /// **'Sent {giftName}. Balance: {diamond}'**
  String liveGiftSent(String giftName, double diamond);
}

class _AppLocalizationsDelegate
    extends LocalizationsDelegate<AppLocalizations> {
  const _AppLocalizationsDelegate();

  @override
  Future<AppLocalizations> load(Locale locale) {
    return SynchronousFuture<AppLocalizations>(lookupAppLocalizations(locale));
  }

  @override
  bool isSupported(Locale locale) =>
      <String>['en', 'zh'].contains(locale.languageCode);

  @override
  bool shouldReload(_AppLocalizationsDelegate old) => false;
}

AppLocalizations lookupAppLocalizations(Locale locale) {
  // Lookup logic when only language code is specified.
  switch (locale.languageCode) {
    case 'en':
      return AppLocalizationsEn();
    case 'zh':
      return AppLocalizationsZh();
  }

  throw FlutterError(
      'AppLocalizations.delegate failed to load unsupported locale "$locale". This is likely '
      'an issue with the localizations generation tool. Please file an issue '
      'on GitHub with a reproducible sample app and the gen-l10n configuration '
      'that was used.');
}
