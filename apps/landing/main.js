document.addEventListener('DOMContentLoaded', () => {
  // Reveal elements on scroll
  const revealObserver = new IntersectionObserver(
    (entries) => {
      entries.forEach((entry) => {
        if (entry.isIntersecting) {
          entry.target.classList.add('visible');
          revealObserver.unobserve(entry.target);
        }
      });
    },
    { threshold: 0.1, rootMargin: '0px 0px -40px 0px' }
  );

  document.querySelectorAll('.reveal').forEach((el) => {
    revealObserver.observe(el);
  });

  // Waitlist form handler
  const form = document.getElementById('waitlist-form');
  const message = document.getElementById('waitlist-message');

  if (form && message) {
    form.addEventListener('submit', (event) => {
      event.preventDefault();
      const email = form.email.value.trim();
      if (!email) return;

      // Placeholder: wire this to your waitlist backend / email service
      message.textContent = `Thanks, ${email}. We'll be in touch.`;
      form.reset();
    });
  }

  // Fetch live GitHub stats (best-effort)
  const repo = 'narcilee7/still';
  const stats = document.querySelectorAll('.stat-value[data-count]');

  async function fetchStats() {
    try {
      const response = await fetch(`https://api.github.com/repos/${repo}`);
      if (!response.ok) return;
      const data = await response.json();

      const mapping = [
        { label: 'Stars', value: data.stargazers_count },
        { label: 'Forks', value: data.forks_count },
        { label: 'Contributors', value: null },
      ];

      stats.forEach((el) => {
        const label = el.nextElementSibling?.textContent;
        const match = mapping.find((m) => m.label === label);
        if (match && typeof match.value === 'number') {
          animateCount(el, match.value);
        }
      });

      // Contributors endpoint is separate
      const contribResponse = await fetch(
        `https://api.github.com/repos/${repo}/contributors?per_page=100`
      );
      if (contribResponse.ok) {
        const contributors = await contribResponse.json();
        const contribEl = Array.from(stats).find(
          (el) => el.nextElementSibling?.textContent === 'Contributors'
        );
        if (contribEl && Array.isArray(contributors)) {
          animateCount(contribEl, contributors.length);
        }
      }
    } catch {
      // Stats are non-critical; leave placeholders on failure
    }
  }

  function animateCount(element, target) {
    const duration = 800;
    const start = performance.now();
    const startValue = 0;

    function step(now) {
      const progress = Math.min((now - start) / duration, 1);
      const value = Math.floor(startValue + (target - startValue) * easeOutCubic(progress));
      element.textContent = value.toLocaleString();
      if (progress < 1) {
        requestAnimationFrame(step);
      }
    }

    requestAnimationFrame(step);
  }

  function easeOutCubic(x) {
    return 1 - Math.pow(1 - x, 3);
  }

  fetchStats();
});
