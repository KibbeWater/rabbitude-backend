from http.server import BaseHTTPRequestHandler, HTTPServer
import json
import whisper

class WhisperHandler(BaseHTTPRequestHandler):
    def do_POST(self):
        global whisper_model
        if self.path == '/api/whisper':
            content_length = int(self.headers['Content-Length'])
            post_data = self.rfile.read(content_length)
            data = json.loads(post_data.decode('utf-8'))
            
            # Process the path from the received JSON
            received_path = data['path']

            print(data)

            # Check if the model is loaded
            if whisper_model is None:
                print("Model not loaded.")
                print(whisper_model)
                self.send_response(400)
                self.send_header('Content-type', 'application/json')
                self.end_headers()
                response = {
                    'status': 'error',
                    'message': 'Model not loaded.'
                }
                self.wfile.write(json.dumps(response).encode('utf-8'))
                return

            # Run the model on the received path
            res = whisper_model.transcribe(received_path)

            # Print the resulting text in red color in the console then reset the color
            print(res['text'])
            print(f"\033[91m{res['text']}\033[0m")
            
            # Send a response back to the client
            self.send_response(200)
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            response = {
                'status': 'success',
                'message': res["text"]
            }
            self.wfile.write(json.dumps(response).encode('utf-8'))
        elif self.path == '/api/init':
            content_length = int(self.headers['Content-Length'])
            post_data = self.rfile.read(content_length)
            data = json.loads(post_data.decode('utf-8'))

            whisper_model = whisper.load_model(data['model'], None, None, True)
            print(f"Model loaded: {whisper_model}")

            # Send a response back to the client
            self.send_response(200)
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            response = {
                'status': 'success',
                'message': f'Model {whisper_model} loaded successfully.'
            }
            self.wfile.write(json.dumps(response).encode('utf-8'))
        else:
            self.send_response(404)
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            response = {
                'status': 'error',
                'message': 'Invalid endpoint.'
            }
            self.wfile.write(json.dumps(response).encode('utf-8'))

def run(server_class=HTTPServer, handler_class=WhisperHandler, port=8118):
    server_address = ('', port)
    httpd = server_class(server_address, handler_class)
    print(f'Starting httpd on port {port}...')
    httpd.serve_forever()

if __name__ == "__main__":
    run()