(function () {
  const state = {
    client: null,
    localAudio: null,
    localVideo: null,
  };

  async function leave() {
    try {
      if (state.localAudio) {
        state.localAudio.close();
        state.localAudio = null;
      }
      if (state.localVideo) {
        state.localVideo.close();
        state.localVideo = null;
      }
      if (state.client) {
        await state.client.leave();
        state.client = null;
      }
    } catch (e) {
      console.warn('saraAgoraRtc.leave', e);
    }
  }

  window.saraAgoraRtc = {
    leave,
    async join(options) {
      const appId = options.appId;
      const channel = options.channel;
      const token = options.token;
      const uid = options.uid;
      const userAccount = (options.userAccount || '').trim();
      const localViewId = options.localViewId;
      const remoteViewId = options.remoteViewId;

      if (!window.AgoraRTC) {
        throw new Error('Agora Web SDK not loaded');
      }

      await leave();

      const client = AgoraRTC.createClient({ mode: 'rtc', codec: 'vp8' });
      state.client = client;

      client.on('user-published', async (user, mediaType) => {
        try {
          await client.subscribe(user, mediaType);
          if (mediaType === 'video' && remoteViewId) {
            const container = document.getElementById(remoteViewId);
            if (container) {
              container.replaceChildren();
              user.videoTrack.play(container);
            }
            if (typeof window.saraAgoraRtcOnRemoteVideo === 'function') {
              window.saraAgoraRtcOnRemoteVideo();
            }
          }
        } catch (e) {
          console.warn('saraAgoraRtc subscribe', e);
        }
      });

      // Token 为 UserAccount 签发时须用字符串进房；仅用数字 uid 且 >4294967295 会 INVALID_PARAMS。
      let joinUid = null;
      if (userAccount.length > 0) {
        joinUid = userAccount;
      } else if (typeof uid === 'number' && uid > 0 && uid <= 4294967295) {
        joinUid = uid;
      } else if (typeof uid === 'string' && uid.length > 0) {
        joinUid = uid;
      }
      if (joinUid == null) {
        throw new Error('Agora uid/userAccount empty or invalid');
      }

      await client.join(appId, channel, token, joinUid);

      const audience = options.role === 'audience';
      if (audience) {
        return true;
      }

      const tracks = await AgoraRTC.createMicrophoneAndCameraTracks();
      state.localAudio = tracks[0];
      state.localVideo = tracks[1];

      if (localViewId) {
        const container = document.getElementById(localViewId);
        if (container) {
          container.replaceChildren();
          state.localVideo.play(container);
        }
      }

      await client.publish(tracks);
      return true;
    },
  };
})();
