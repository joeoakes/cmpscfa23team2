from flask import Flask, jsonify, request

app = Flask(__name__)

# In-memory storage of items
items = []

@app.route('/')
def home():
    return jsonify(message="Hello, World!")

@app.route('/api/items', methods=['GET'])
def get_all_items():
    return jsonify(items=items)

@app.route('/api/items', methods=['POST'])
def create_item():
    data = request.json
    item = {
        'id': len(items) + 1,
        'name': data['name']
    }
    items.append(item)
    return jsonify(item=item), 201

@app.route('/api/items/<int:item_id>', methods=['GET'])
def get_item(item_id):
    item = next((item for item in items if item['id'] == item_id), None)
    if item:
        return jsonify(item=item)
    return jsonify(message="Item not found"), 404

@app.route('/api/items/<int:item_id>', methods=['PUT'])
def update_item(item_id):
    data = request.json
    item = next((item for item in items if item['id'] == item_id), None)
    if item:
        item['name'] = data['name']
        return jsonify(item=item)
    return jsonify(message="Item not found"), 404

@app.route('/api/items/<int:item_id>', methods=['DELETE'])
def delete_item(item_id):
    global items
    items = [item for item in items if item['id'] != item_id]
    return jsonify(message="Item deleted"), 200

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5005)


