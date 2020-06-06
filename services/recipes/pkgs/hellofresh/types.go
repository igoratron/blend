package hellofresh

type Id string

type Recipe struct {
	Id          string       `dynamodbav:"id"`
	Name        string       `dynamodbav:"name"`
	WebsiteUrl  string       `dynamodbav:"websiteUrl"`
	Ingredients []Ingredient `dynamodbav:"ingredients"`
	ImagePath   string       `dynamodbav:"imagePath"`
	Yields      []Yield      `dynamodbav:"yields"`
}

type Yield struct {
	Yields      uint              `dynamodbav:"yields"`
	Ingredients []YieldIngredient `dynamodbav:"ingredients"`
}

type YieldIngredient struct {
	Id     string  `dynamodbav:"id"`
	Amount float32 `dynamodbav:"amount"`
	Unit   string  `dynamodbav:"unit"`
}

type Ingredient struct {
	Id   string `dynamodbav:"id"`
	Name string `dynamodbav:"name"`
}
