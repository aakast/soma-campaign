// Navigation Toggle
document.addEventListener('DOMContentLoaded', function() {
    const navToggle = document.getElementById('navToggle');
    const navMenu = document.getElementById('navMenu');
    const navbar = document.getElementById('navbar');
    
    // Mobile menu toggle
    if (navToggle) {
        navToggle.addEventListener('click', function() {
            navMenu.classList.toggle('active');
            
            // Animate hamburger menu
            const spans = navToggle.querySelectorAll('span');
            if (navMenu.classList.contains('active')) {
                spans[0].style.transform = 'rotate(45deg) translateY(8px)';
                spans[1].style.opacity = '0';
                spans[2].style.transform = 'rotate(-45deg) translateY(-8px)';
            } else {
                spans[0].style.transform = 'none';
                spans[1].style.opacity = '1';
                spans[2].style.transform = 'none';
            }
        });
    }
    
    // Close mobile menu when clicking a link
    const navLinks = document.querySelectorAll('.nav-link');
    navLinks.forEach(link => {
        link.addEventListener('click', () => {
            navMenu.classList.remove('active');
            const spans = navToggle.querySelectorAll('span');
            spans[0].style.transform = 'none';
            spans[1].style.opacity = '1';
            spans[2].style.transform = 'none';
        });
    });
    
    // Navbar scroll effect
    window.addEventListener('scroll', function() {
        if (window.scrollY > 50) {
            navbar.classList.add('scrolled');
        } else {
            navbar.classList.remove('scrolled');
        }
    });
    
    // Hero video scroll behavior
    const heroVideo = document.getElementById('heroVideo');
    const heroVideoElement = document.getElementById('heroVideoElement');
    const mainContent = document.getElementById('mainContent');
    
    if (heroVideo && heroVideoElement) {
        let isVideoPlaying = true;
        
        // Stop video when scrolling down
        window.addEventListener('scroll', function() {
            const scrollPosition = window.scrollY;
            const videoHeight = heroVideo.offsetHeight;
            
            if (scrollPosition > 100 && isVideoPlaying) {
                heroVideoElement.pause();
                isVideoPlaying = false;
                
                // Fade out video section
                heroVideo.style.opacity = Math.max(0, 1 - (scrollPosition / videoHeight));
            } else if (scrollPosition <= 100 && !isVideoPlaying) {
                heroVideoElement.play();
                isVideoPlaying = true;
                heroVideo.style.opacity = 1;
            }
        });
        
        // Smooth scroll on arrow click
        const scrollIndicator = document.querySelector('.hero-scroll-indicator');
        if (scrollIndicator) {
            scrollIndicator.addEventListener('click', function() {
                mainContent.scrollIntoView({ behavior: 'smooth' });
            });
        }
    }
    
    // Add playful hover effects to doodle elements
    const doodleFrames = document.querySelectorAll('.doodle-frame, .doodle-frame-small, .doodle-circle');
    doodleFrames.forEach(frame => {
        frame.addEventListener('mouseenter', function() {
            this.style.transform = `rotate(${Math.random() * 10 - 5}deg) scale(1.05)`;
        });
        
        frame.addEventListener('mouseleave', function() {
            this.style.transform = 'rotate(-2deg) scale(1)';
        });
    });
    
    // Lazy loading for images
    const images = document.querySelectorAll('img[data-src]');
    const imageOptions = {
        threshold: 0,
        rootMargin: '0px 0px 50px 0px'
    };
    
    const imageObserver = new IntersectionObserver((entries, observer) => {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                const img = entry.target;
                img.src = img.dataset.src;
                img.classList.add('loaded');
                observer.unobserve(img);
            }
        });
    }, imageOptions);
    
    images.forEach(img => imageObserver.observe(img));
    
	// Add animation to elements when they come into view
	const animateElements = document.querySelectorAll('.issue-card, .news-card, .intro-wrapper, .fb-card');
    const animateOptions = {
        threshold: 0.1,
        rootMargin: '0px 0px -50px 0px'
    };
    
    const animateObserver = new IntersectionObserver((entries) => {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                entry.target.style.animation = 'fadeInUp 0.6s ease forwards';
            }
        });
    }, animateOptions);
    
	animateElements.forEach(el => animateObserver.observe(el));

	// Facebook feed hydration
	const fbContainer = document.getElementById('facebook-feed');
	if (fbContainer) {
		fetch('/api/facebook/feed')
			.then(res => res.json())
			.then(data => {
				if (!data || !Array.isArray(data.posts) || data.posts.length === 0) {
					return;
				}

				// Hide plugin embed if we have posts (CSS also hides via has-posts class)
				const fbPlugin = document.querySelector('.facebook-wrapper .fb-page');
				if (fbPlugin) fbPlugin.style.display = 'none';

				fbContainer.classList.add('has-posts');
				const fragment = document.createDocumentFragment();

				const formatDate = (iso) => {
					try {
						const d = new Date(iso);
						return d.toLocaleDateString('da-DK', { day: 'numeric', month: 'long', year: 'numeric' });
					} catch (_) {
						return '';
					}
				};

				const truncate = (text, n) => {
					if (!text) return '';
					return text.length > n ? text.slice(0, n - 1) + '…' : text;
				};

				data.posts.forEach(post => {
					const card = document.createElement('article');
					card.className = 'fb-card';

					if (post.full_picture) {
						const imgWrap = document.createElement('div');
						imgWrap.className = 'fb-card-image';
						const img = document.createElement('img');
						img.src = post.full_picture;
						img.alt = 'Facebook opslag';
						imgWrap.appendChild(img);
						card.appendChild(imgWrap);
					}

					const content = document.createElement('div');
					content.className = 'fb-card-content';

					const date = document.createElement('time');
					date.className = 'fb-card-date';
					date.textContent = formatDate(post.created_time);
					content.appendChild(date);

					if (post.message) {
						const p = document.createElement('p');
						p.className = 'fb-card-message';
						p.textContent = truncate(post.message, 220);
						content.appendChild(p);
					}

					if (post.permalink_url) {
						const a = document.createElement('a');
						a.className = 'fb-card-link';
						a.href = post.permalink_url;
						a.target = '_blank';
						a.rel = 'noopener noreferrer';
						a.textContent = 'Læs på Facebook →';
						content.appendChild(a);
					}

					card.appendChild(content);
					fragment.appendChild(card);
				});

				fbContainer.appendChild(fragment);
				fbContainer.querySelectorAll('.fb-card').forEach(el => animateObserver.observe(el));
			})
			.catch(() => { /* silent */ });
	}
});

// Add fadeInUp animation
const style = document.createElement('style');
style.textContent = `
    @keyframes fadeInUp {
        from {
            opacity: 0;
            transform: translateY(30px);
        }
        to {
            opacity: 1;
            transform: translateY(0);
        }
    }
`;
document.head.appendChild(style);

// TinaCMS Integration
if (typeof TinaCMS !== 'undefined') {
    // Initialize TinaCMS
    const cms = new TinaCMS({
        enabled: true,
        sidebar: true,
        toolbar: true,
    });
    
    // Add forms for editable content
    cms.registerForm({
        id: 'home-hero',
        label: 'Hero Section',
        fields: [
            {
                name: 'title',
                label: 'Title',
                component: 'text',
            },
            {
                name: 'subtitle',
                label: 'Subtitle',
                component: 'text',
            },
            {
                name: 'videoUrl',
                label: 'Video URL',
                component: 'text',
            },
        ],
        onSubmit: async (values) => {
            // Send to backend
            await fetch('/api/tina/content', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    collection: 'hero',
                    id: 'home',
                    data: values,
                }),
            });
        },
    });
}