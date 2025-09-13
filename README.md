# Soma Mayel Campaign Website

A modern, responsive campaign website for Soma Mayel, lead candidate (Spidskandidat) for Radikale Venstre in Fredensborg Kommune.

## Features

- 🎨 **Beautiful Design**: Incorporates Radikale Venstre's brand colors with a playful "children's drawings" aesthetic
- 🎥 **Fullscreen Video Hero**: Engaging landing page with autoplay video that pauses on scroll
- 📱 **Fully Responsive**: Works perfectly on all devices
- ✏️ **CMS Integration**: Built-in admin UI with Basic Auth and file-based JSON content
- 📘 **Facebook Integration**: Embedded Facebook feed for social media engagement
- 🚀 **Fast & Lightweight**: Built with Go and Fiber framework for optimal performance
- 🐳 **Docker Ready**: Easy deployment with Docker and Docker Compose

## Tech Stack

- **Backend**: Go with Fiber web framework
- **Templates**: HTML with Go's template engine
- **Styling**: Custom CSS with Radikale Venstre color scheme
- **CMS**: Built-in admin UI (Tina-like) for content management
- **Deployment**: Docker & Docker Compose

## Quick Start

### Prerequisites

- Docker and Docker Compose installed
- OR Go 1.21+ installed locally

### Using Docker (Recommended)

1. Clone the repository:
```bash
git clone <repository-url>
cd <repository-name>
```

2. Copy the environment file and configure:
```bash
cp .env.example .env
# Edit .env with your settings
```

3. Build and run with Docker Compose:
```bash
docker-compose up --build
```

4. Open your browser and visit: `http://localhost:3000`

### Local Development

1. Install dependencies:
```bash
go mod download
```

2. Copy and configure environment:
```bash
cp .env.example .env
# Edit .env with your settings
```

3. Run the application:
```bash
go run main.go
```

4. Visit: `http://localhost:3000`

## Project Structure

```
├── main.go                 # Application entry point
├── handlers/               # HTTP request handlers
│   ├── home.go
│   ├── about.go
│   ├── politics.go
│   ├── news.go
│   ├── contact.go
│   └── tina.go            # Legacy TinaCMS API handlers (optional)
├── models/                 # Data models
│   ├── post.go
│   └── content.go
├── templates/              # HTML templates
│   ├── layouts/
│   │   └── main.html      # Main layout
│   ├── home.html
│   ├── about.html
│   ├── politics.html
│   ├── news.html
│   ├── contact.html
│   ├── blog-post.html
│   └── 404.html
├── static/                 # Static assets
│   ├── css/
│   │   └── main.css       # Main stylesheet
│   ├── js/
│   │   └── main.js        # JavaScript interactions
│   ├── images/            # Image assets
│   └── videos/            # Video files
├── content/                # CMS content (JSON files)
│   ├── posts/             # Blog posts
│   └── pages/             # Static pages
├── templates/
│   └── admin.html         # Admin UI
├── Dockerfile
├── docker-compose.yml
└── .env.example
```

## Configuration

### Environment Variables

Create a `.env` file based on `.env.example`:

- `PORT`: Application port (default: 3000)
- `ENV`: Environment (development/production)
- `ADMIN_USERNAME`: Basic auth username for admin (default: admin)
- `ADMIN_PASSWORD`: Basic auth password for admin (default: admin123)
- `FACEBOOK_PAGE_ID`: Facebook page for social feed
- `CONTACT_EMAIL`: Email for contact form submissions

### Admin CMS

Access the admin UI at `/admin` (protected by Basic Auth). From here you can:
- Create, edit, delete blog posts
- Upload images (stored under `static/images/uploads/`)

Content is stored as JSON under `content/posts/`.

### Adding Content

#### Blog Posts
Create JSON files in `content/posts/` with the following structure:
```json
{
  "id": "unique-id",
  "title": "Post Title",
  "slug": "post-slug",
  "content": "Post content...",
  "excerpt": "Short description",
  "author": "Soma Mayel",
  "date": "2025-01-15T10:00:00Z",
  "image": "/static/images/post-image.jpg",
  "tags": ["Politik", "Fredensborg"],
  "isFeatured": true
}
```

#### Adding Videos
Place video files in `static/videos/` and reference them in templates or content.

#### Adding Images
Place images in `static/images/` for use throughout the site.

## Deployment

### Using Docker on a VM

1. Copy the project to your VM
2. Ensure Docker and Docker Compose are installed
3. Configure your `.env` file
4. Run:
```bash
docker-compose up -d
```

### Using Systemd (Alternative)

1. Build the Go binary:
```bash
go build -o soma-campaign main.go
```

2. Create a systemd service file:
```bash
sudo nano /etc/systemd/system/soma-campaign.service
```

3. Add the service configuration:
```ini
[Unit]
Description=Soma Mayel Campaign Website
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/path/to/repository
ExecStart=/path/to/repository/soma-campaign
Restart=on-failure
Environment=PORT=3000

[Install]
WantedBy=multi-user.target
```

4. Enable and start the service:
```bash
sudo systemctl enable soma-campaign
sudo systemctl start soma-campaign
```

## Maintenance

### Updating Content
- Use TinaCMS interface at `/admin`
- Or directly edit JSON files in `content/` directory

### Backup
Regular backups should include:
- `content/` directory (all CMS content)
- `static/images/` and `static/videos/` (media files)
- `.env` file (configuration)

### Monitoring
- Check application logs: `docker-compose logs -f web`
- Monitor system resources
- Set up health checks for the `/` endpoint

## Security Considerations

- Always use HTTPS in production (configure reverse proxy)
- Keep dependencies updated
- Secure your `.env` file (never commit to git)
- Configure proper CORS settings for production
- Implement rate limiting for contact form
- Regular security updates

## Support

For issues or questions about the website, please contact the development team.

## License

© 2025 Soma Mayel Campaign. All rights reserved.