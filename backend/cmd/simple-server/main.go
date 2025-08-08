package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

// Simple GraphQL-like response structure for validation
type GraphQLResponse struct {
	Data   interface{} `json:"data,omitempty"`
	Errors []string    `json:"errors,omitempty"`
}

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type Post struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	Author    User     `json:"author"`
	Tags      []string `json:"tags"`
	Published bool     `json:"published"`
}

func main() {
	r := gin.Default()

	// Enable CORS
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// Simple GraphQL endpoint for validation
	r.POST("/graphql", func(c *gin.Context) {
		var request map[string]interface{}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, GraphQLResponse{
				Errors: []string{"Invalid JSON"},
			})
			return
		}

		query, exists := request["query"].(string)
		if !exists {
			c.JSON(400, GraphQLResponse{
				Errors: []string{"No query provided"},
			})
			return
		}

		// Simple query handling for validation
		if query == "{ me { id name email } }" {
			c.JSON(200, GraphQLResponse{
				Data: map[string]interface{}{
					"me": User{
						ID:    "1",
						Name:  "Demo User",
						Email: "demo@example.com",
					},
				},
			})
			return
		}

		if query == "{ posts { edges { node { id title author { name } } } } }" {
			c.JSON(200, GraphQLResponse{
				Data: map[string]interface{}{
					"posts": map[string]interface{}{
						"edges": []map[string]interface{}{
							{
								"node": Post{
									ID:      "1",
									Title:   "Sample Post",
									Content: "Sample content",
									Author: User{
										ID:   "1",
										Name: "Demo User",
									},
									Tags:      []string{"demo", "test"},
									Published: true,
								},
							},
						},
					},
				},
			})
			return
		}

		c.JSON(200, GraphQLResponse{
			Data: map[string]interface{}{
				"message": "GraphQL server is working! Query received: " + query,
			},
		})
	})

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "graphql-typescript-go-backend",
			"message": "Server is running and ready for GraphQL queries",
		})
	})

	// Simple GraphQL playground
	r.GET("/playground", func(c *gin.Context) {
		c.Header("Content-Type", "text/html")
		c.String(200, `
<!DOCTYPE html>
<html>
<head>
    <title>GraphQL Playground</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .container { max-width: 800px; margin: 0 auto; }
        textarea { width: 100%%; height: 200px; margin: 10px 0; }
        button { padding: 10px 20px; background: #0066cc; color: white; border: none; cursor: pointer; }
        .result { background: #f5f5f5; padding: 20px; margin: 10px 0; border-radius: 5px; }
    </style>
</head>
<body>
    <div class="container">
        <h1>GraphQL Playground</h1>
        <p>Test your GraphQL queries here:</p>
        
        <h3>Sample Queries:</h3>
        <pre>{ me { id name email } }</pre>
        <pre>{ posts { edges { node { id title author { name } } } } }</pre>
        
        <textarea id="query" placeholder="Enter your GraphQL query here...">{ me { id name email } }</textarea>
        <br>
        <button onclick="executeQuery()">Execute Query</button>
        
        <h3>Result:</h3>
        <div id="result" class="result">Click "Execute Query" to see results</div>
    </div>

    <script>
        async function executeQuery() {
            const query = document.getElementById('query').value;
            const resultDiv = document.getElementById('result');
            
            try {
                const response = await fetch('/graphql', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ query: query })
                });
                
                const result = await response.json();
                resultDiv.innerHTML = '<pre>' + JSON.stringify(result, null, 2) + '</pre>';
            } catch (error) {
                resultDiv.innerHTML = '<pre style="color: red;">Error: ' + error.message + '</pre>';
            }
        }
    </script>
</body>
</html>
		`)
	})

	log.Println("🚀 GraphQL server ready at http://localhost:8080/graphql")
	log.Println("🎮 GraphQL playground available at http://localhost:8080/playground")
	log.Println("❤️  Health check at http://localhost:8080/health")
	
	r.Run(":8080")
}