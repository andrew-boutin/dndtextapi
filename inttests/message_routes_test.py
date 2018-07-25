# Copyright (C) 2018, Baking Bits Studios - All Rights Reserved

import requests
from base import BaseTest

class MessageRoutesTest(BaseTest):

    def setUp(self):
        super(MessageRoutesTest, self).setUp()
        self.url = f"{self.base}/messages"

    def tearDown(self):
        super(MessageRoutesTest, self).tearDown()