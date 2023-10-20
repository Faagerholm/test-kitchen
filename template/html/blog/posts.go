package blog

import (
	"html/template"
)

type Post struct {
	ID      int
	Title   string
	Tags    []string
	Content template.HTML
}

func GetPosts() []Post {
	return []Post{
		{
			ID:      1,
			Title:   "Introduction to Go",
			Tags:    []string{"Go", "Programming"},
			Content: template.HTML("Go is a statically typed, compiled language."),
		},
		{
			ID:      2,
			Title:   "Web Development with Go",
			Tags:    []string{"Go", "Web"},
			Content: template.HTML("Go is a great choice for building web applications."),
		},
		{
			ID:      3,
			Title:   "Working with Databases in Go",
			Tags:    []string{"Go", "Database"},
			Content: template.HTML("Go has excellent support for various databases."),
		},
		{
			ID:      4,
			Title:   "Concurrency in Go",
			Tags:    []string{"Go", "Concurrency"},
			Content: template.HTML("Go makes it easy to write concurrent programs."),
		},
		{
			ID:      5,
			Title:   "Error Handling in Go",
			Tags:    []string{"Go", "Error"},
			Content: template.HTML("Go provides robust error handling mechanisms."),
		},
		{
			ID:      6,
			Title:   "Working with Files in Go",
			Tags:    []string{"Go", "File"},
			Content: template.HTML("Go has powerful file handling capabilities."),
		},
		{
			ID:      7,
			Title:   "Creating REST APIs with Go",
			Tags:    []string{"Go", "API"},
			Content: template.HTML("Go is a popular choice for building RESTful APIs."),
		},
		{
			ID:      8,
			Title:   "Testing in Go",
			Tags:    []string{"Go", "Testing"},
			Content: template.HTML("Go has a built-in testing framework."),
		},
		{
			ID:      9,
			Title:   "Deploying Go Applications",
			Tags:    []string{"Go", "Deployment"},
			Content: template.HTML("Go applications are easy to deploy and scale."),
		},
		{
			ID:      10,
			Title:   "Go Conventions and Best Practices",
			Tags:    []string{"Go", "Best Practices"},
			Content: template.HTML("Follow Go's conventions for clean and readable code."),
		},
		{
			ID:      11,
			Title:   "Getting Started with HTMX",
			Tags:    []string{"HTMX", "Web Development"},
			Content: template.HTML("HTMX is a library for creating modern, dynamic web applications with minimal JavaScript."),
		},
		{
			ID:      12,
			Title:   "HTMX vs. AJAX",
			Tags:    []string{"HTMX", "AJAX"},
			Content: template.HTML("Learn the differences between HTMX and traditional AJAX for enhancing web interactions."),
		},
		{
			ID:      13,
			Title:   "Using HTMX with Go",
			Tags:    []string{"HTMX", "Go", "Web"},
			Content: template.HTML("Integrate HTMX seamlessly into your Go web applications."),
		},
		{
			ID:      14,
			Title:   "HTMX and Server-Sent Events (SSE)",
			Tags:    []string{"HTMX", "SSE", "Real-time"},
			Content: template.HTML("Explore real-time web updates using HTMX and Server-Sent Events."),
		},
		{
			ID:      15,
			Title:   "Building a Shopping Cart with HTMX",
			Tags:    []string{"HTMX", "Web Development", "Shopping Cart"},
			Content: template.HTML("Create a dynamic shopping cart using HTMX for a smooth user experience."),
		},
		{
			ID:      16,
			Title:   "HTMX and Progressive Enhancement",
			Tags:    []string{"HTMX", "Progressive Enhancement"},
			Content: template.HTML("Discover how HTMX can help implement progressive enhancement in web applications."),
		},
		{
			ID:      17,
			Title:   "Authentication with HTMX",
			Tags:    []string{"HTMX", "Authentication", "Security"},
			Content: template.HTML("Secure your HTMX-powered applications with proper authentication mechanisms."),
		},
		{
			ID:      18,
			Title:   "HTMX and SEO",
			Tags:    []string{"HTMX", "SEO", "Search Engine Optimization"},
			Content: template.HTML("Learn about SEO considerations when using HTMX for client-side rendering."),
		},
		{
			ID:      19,
			Title:   "HTMX Best Practices",
			Tags:    []string{"HTMX", "Best Practices", "Web Development"},
			Content: template.HTML("Follow best practices for efficient and maintainable HTMX development."),
		},
		{
			ID:      20,
			Title:   "Migrating from jQuery to HTMX",
			Tags:    []string{"HTMX", "jQuery", "Migration"},
			Content: template.HTML("Tips and strategies for migrating legacy jQuery code to HTMX."),
		},
	}
}
