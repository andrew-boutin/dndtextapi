import unittest, requests

class BaseTest(unittest.TestCase):

    def __init__(self, *args, **kwargs):
        super(BaseTest, self).__init__(*args, **kwargs)
        self.base = "http://app:8080/"

    def setUp(self):
        pass

    def tearDown(self):
        pass

    def test_get_channels(self):
        """Test the public get channels route."""
        # Make sure not including the accept header isn't allowed
        r = requests.get(f"{self.base}public/channels")
        self.assertEqual(r.status_code, 400)

        # Add the header and make sure the call is ok
        headers = {"Accept": "application/json"}
        r = requests.get(f"{self.base}public/channels", headers=headers)
        self.assertEqual(r.status_code, 200)
        # TODO: assert the contents

    def test_get_channel(self):
        """Test the public get channel route."""
        # Make sure not including the accept header isn't allowed
        url = f"{self.base}public/channels/"
        r = requests.get(f"{url}1")
        self.assertEqual(r.status_code, 400)

        # Add the header and make sure the call is ok
        headers = {"Accept": "application/json"}
        r = requests.get(f"{url}1", headers=headers)
        self.assertEqual(r.status_code, 200)
        # TODO: assert the contents

        # Verify we get a not found on a made up id
        r = requests.get(f"{url}999", headers=headers)
        self.assertEqual(r.status_code, 404)

        # Make sure we can't get a private channel
        url = f"{url}2"
        r = requests.get(url, headers=headers)
        self.assertEqual(r.status_code, 401)

    def test_get_messages_from_public_channel(self):
        """Test the public get messages route for getting messages from a channel."""
        # Make sure not including the accept header isn't allowed
        url = f"{self.base}public/messages"
        r = requests.get(url)
        self.assertEqual(r.status_code, 400)

        # Verify 400 when not including the required query param
        headers = {"Accept": "application/json"}
        r = requests.get(url, headers=headers)
        self.assertEqual(r.status_code, 400)

        url = f"{url}?channelID="

        # Use a private channel id verify denied
        r = requests.get(f"{url}2", headers=headers)
        self.assertEqual(r.status_code, 401)

        # Make up a channel id verify not found
        r = requests.get(f"{url}999", headers=headers)
        self.assertEqual(r.status_code, 404)

        # Use a public channel for a valid request
        r = requests.get(f"{url}1", headers=headers)
        self.assertEqual(r.status_code, 200)
        # TODO: assert the contents - no meta messages