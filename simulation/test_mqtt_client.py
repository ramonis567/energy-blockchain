import unittest
from unittest.mock import patch, MagicMock
import sys
import os

# Add the src directory to the Python path
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), 'src')))

from mqtt_client import connect, disconnect, publish_data

class TestMqttClient(unittest.TestCase):

    @patch('mqtt_client.client')
    def test_connect_publish_disconnect(self, mock_client_instance):
        # Test connect
        connect()
        mock_client_instance.connect.assert_called_once()

        # Test publish
        topic = "test/topic"
        data = "test_data"
        publish_data(topic, data)
        mock_client_instance.publish.assert_called_once_with(topic, data)

        # Test disconnect
        disconnect()
        mock_client_instance.disconnect.assert_called_once()

if __name__ == '__main__':
    unittest.main()
