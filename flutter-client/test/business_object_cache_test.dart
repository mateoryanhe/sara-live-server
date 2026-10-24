import 'dart:io';

import 'package:flutter_test/flutter_test.dart';
import 'package:sara_live_app/core/cache/api_disk_cache.dart';
import 'package:sara_live_app/core/cache/business_object_cache.dart';
import 'package:sara_live_app/core/cache/device_tier.dart';
import 'package:sara_live_app/core/cache/memory_cache_store.dart';
import 'package:sara_live_app/data/user/user_profile.dart';

void main() {
  test('L2 disk restores into L1 after memory clear', () async {
    final temp = await Directory.systemTemp.createTemp('api_cache_test_');
    addTearDown(() async {
      if (temp.existsSync()) {
        await temp.delete(recursive: true);
      }
      ApiDiskCache.instance.dispose();
    });

    await ApiDiskCache.init(DeviceTier.mid, rootOverride: temp);
    final memory = MemoryCacheStore(tier: DeviceTier.mid);
    final cache = BusinessObjectCache(
      memory: memory,
      disk: ApiDiskCache.instance,
    );

    const userId = 10001;
    final profile = UserProfile(
      userId: '10001',
      nickname: 'n',
      avatar: 'a',
      gold: 1,
      diamond: 2,
      shareCode: 's',
      vipLevel: 0,
      isAnchor: false,
      prettyId: '',
      gender: 0,
      followCount: 0,
      followerCount: 0,
    );

    await cache.putUser(
      userId,
      'profile',
      profile,
      encode: (UserProfile v) => v.toJson(),
    );
    memory.clear();

    final restored = await cache.getUser<UserProfile>(
      userId,
      'profile',
      decode: UserProfile.fromJson,
    );
    expect(restored?.nickname, 'n');
    expect(memory.getUser<UserProfile>(userId, 'profile')?.nickname, 'n');
  });
}
