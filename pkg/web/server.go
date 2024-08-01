package web

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/rusik69/iamrotator/pkg/aws"
	"github.com/rusik69/iamrotator/pkg/config"
	"github.com/rusik69/iamrotator/pkg/types"
	"github.com/sirupsen/logrus"
)

var cfg config.Config

var (
	cachedKeys      []types.AWSAccessKey
	cacheMutex      sync.Mutex
	cacheExpiration time.Time
	cacheDuration   = 5 * time.Minute // Cache duration
)

// Listen starts the web server.
func Listen(configPath string) error {
	cf, err := config.Load(configPath)
	if err != nil {
		return err
	}
	cfg = cf
	http.HandleFunc("/", serveHomePage)
	http.HandleFunc("/listkeys", serveListKeys)
	logrus.Infof("Listening on %s", cfg.Web.ListenAddr)
	http.ListenAndServe(cfg.Web.ListenAddr, nil)
	return nil
}

// serveHomePage serves the home page.
func serveHomePage(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("Serving home page")
	html := `
    <!DOCTYPE html>
    <html>
    <head>
        <title>IAM Rotator</title>
        <style>
            #loadingCircle {
                display: none;
                border: 16px solid #f3f3f3;
                border-top: 16px solid #3498db;
                border-radius: 50%;
                width: 120px;
                height: 120px;
                animation: spin 2s linear infinite;
                position: absolute;
                left: 50%;
                top: 50%;
                transform: translate(-50%, -50%);
            }

            @keyframes spin {
                0% { transform: rotate(0deg); }
                100% { transform: rotate(360deg); }
            }
        </style>
    </head>
    <body>
        <h1>IAM Rotator</h1>
        <button onclick="getKeys()">Get Key List</button>
        <div id="loadingCircle"></div>
        <table id="keyTable" border="1">
            <tr>
                <th>Username</th>
                <th>Access Key ID</th>
                <th>Account ID</th>
                <th>Account Name</th>
                <th>Status</th>
                <th>Create Date</th>
            </tr>
        </table>
        <script>
            function getKeys() {
                showLoadingCircle();
                fetch('/listkeys')
                    .then(response => response.json())
                    .then(data => {
                        const table = document.getElementById('keyTable');
                        data.forEach(key => {
                            const row = table.insertRow();
                            const cell1 = row.insertCell(0);
                            const cell2 = row.insertCell(1);
                            const cell3 = row.insertCell(2);
                            const cell4 = row.insertCell(3);
                            const cell5 = row.insertCell(4);
                            const cell6 = row.insertCell(5);
                            cell1.textContent = key.UserName;
                            cell2.textContent = key.AccessKeyID;
                            cell3.textContent = key.AccountID;
                            cell4.textContent = key.AccountName;
                            cell5.textContent = key.Status;
                            cell6.textContent = new Date(key.CreateDate).toLocaleString();
                         });
                        hideLoadingCircle();
                    })
                    .catch(error => {
                        console.error('Error fetching key list:', error);
                        hideLoadingCircle();
                    });
            }

            function showLoadingCircle() {
                const loadingCircle = document.getElementById('loadingCircle');
                loadingCircle.style.display = 'block';
            }

            function hideLoadingCircle() {
                const loadingCircle = document.getElementById('loadingCircle');
                loadingCircle.style.display = 'none';
            }
        </script>
    </body>
    </html>
    `
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

// serveListKeys serves the list keys page.
func serveListKeys(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("Serving list keys")
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	if time.Now().After(cacheExpiration) {
		sess, err := aws.CreateSession(cfg.AWS)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		keys, err := aws.ListAccessKeys(sess, cfg.AWS)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cachedKeys = keys
		cacheExpiration = time.Now().Add(cacheDuration)
	}
	json.NewEncoder(w).Encode(cachedKeys)
}
