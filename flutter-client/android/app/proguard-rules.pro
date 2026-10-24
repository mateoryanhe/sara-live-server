# Flutter / plugins (Agora RTC)
-keep class io.flutter.** { *; }
-keep class io.agora.** { *; }
-dontwarn io.agora.**

# Optional Play Core (Flutter deferred components; not used in this app)
-dontwarn com.google.android.play.core.**
