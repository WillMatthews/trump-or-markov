# Real Data

This README details how to access the real data used in this project (which is required to train the model,
and therefore run the backend).


>[!IMPORTANT]
> You will need to download the dataset to run the project. The dataset is not included in this repository due to its size.

>[!WARNING]
> If the dataset exists but you still encounter errors, make sure the filename is the same as the one in the config.yaml file in the `cloud` directory.

## Datasets

### Trump Tweets

I used the Trump Archive dataset, which contains a mixture of Truth social and Twitter data.
The dataset can be found [here](https://www.thetrumparchive.com/).

Their website conveniently provides a JSON dataset of all the tweets, which can be downloaded [here](https://drive.google.com/file/d/16wm-2NTKohhcA26w-kaWfhLIGwl_oX95/view).

```bash
cat tweets_01-08-2021.json | jq | grep "thin person drinking Diet Coke" -B 2 -A 8
  {
    "id": 334168974700982300,
    "text": "\"\"@KarltheMarx: “@realDonaldTrump: I have never seen a thin person drinking Diet Coke.”",
    "isRetweet": "f",
    "isDeleted": "f",
    "device": "Twitter for Android",
    "favorites": 158,
    "retweets": 219,
    "date": "2013-05-14 04:51:06",
    "isFlagged": "f"
  },
```


### Moby Dick

Moby Dick is available from [Project Gutenberg](https://www.gutenberg.org/ebooks/2701).
Download the plain text, and clean as appropriate.


## Dilution Datasets *WIP*

### The King James Bible
https://www.gutenberg.org/ebooks/10

### Bible (but only quotes from Jesus Christ of Nazareth)
https://www.eldoradoweather.com/current/jesus-quotes.php
