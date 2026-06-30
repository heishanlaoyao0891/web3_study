/* app.js — 共享前端交互 */

// Toast 通知（全局复用）
function showToast(message, isError) {
  var toast = document.getElementById('appToast') || (function () {
    var t = document.createElement('div');
    t.id = 'appToast';
    t.className = 'app-toast';
    document.body.appendChild(t);
    return t;
  })();
  toast.textContent = message;
  toast.style.background = isError ? '#F43F5E' : '#10B981';
  toast.style.color = '#fff';
  toast.style.display = 'block';
  setTimeout(function () { toast.style.display = 'none'; }, 3000);
}

// 打卡功能（首页调用）
function doCheckin() {
  var btn = document.getElementById('checkinBtn');
  if (!btn) return;
  btn.disabled = true;
  btn.textContent = '打卡中...';
  fetch('/api/checkin', { method: 'POST' })
    .then(function (res) { return res.json(); })
    .then(function (data) {
      if (data.error) {
        showToast(data.error, true);
        btn.textContent = '立即打卡';
        btn.disabled = false;
      } else {
        showToast('打卡成功！经验+' + data.exp_gained + ' 金币+' + data.coins_gained);
        btn.textContent = '今日已打卡';
        setTimeout(function () { location.reload(); }, 1500);
      }
    })
    .catch(function () {
      showToast('打卡失败，请重试', true);
      btn.textContent = '立即打卡';
      btn.disabled = false;
    });
}

// 首页加载时检查打卡状态
document.addEventListener('DOMContentLoaded', function () {
  var btn = document.getElementById('checkinBtn');
  if (btn && btn.dataset.checkinStatus === 'checked') {
    btn.textContent = '今日已打卡';
    btn.disabled = true;
  }
});