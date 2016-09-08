package main

import(
	"log"
	"time"
	"reflect"
	"net/http"
	"strconv"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Topic struct {
	ID int
	Title string
	Description string
	Likes int
	Comments []Comment
	Created time.Time
}

type ID struct {
	ID int
}

type Comment struct {
	Author string
	Comment string
	Date time.Time
}

type Disc struct {
	Title string
	Description string
	Comments []Comment
	Date time.Time
}

var lastId int
func getItem(id string) Topic {
	x, _ := strconv.Atoi(id)
	sess := session.New(&aws.Config{Region: aws.String("us-west-2")})

	svc := dynamodb.New(sess)
	
	d := ID{ID: x}

	av, err := dynamodbattribute.MarshalMap(d)
	if err != nil {
		log.Println("Error converting ID to Attribute Value")
		log.Println(err)
	}

	params := &dynamodb.GetItemInput{TableName: aws.String("discussions"), Key: av, AttributesToGet: []*string{aws.String("ID"), aws.String("Title"), aws.String("Description"), aws.String("Likes"), aws.String("Comments"), aws.String("Created")}, ConsistentRead: aws.Bool(true),}

	resp, err := svc.GetItem(params)
	if err != nil {
		log.Println("Error Getting Item.")
		log.Println(err)
	}

	m := Topic{}
	dynamodbattribute.UnmarshalMap(resp.Item, &m)
	
	return m
}

func getLastItemID() int {
	return lastId
}

func getTableItems() []Topic {
	sess := session.New(&aws.Config{Region: aws.String("us-west-2")})

	svc := dynamodb.New(sess)

	params := &dynamodb.ScanInput{TableName: aws.String("discussions"), AttributesToGet: []*string{aws.String("ID"), aws.String("Title"), aws.String("Description"), aws.String("Likes"), aws.String("Comments"), aws.String("Created")}}

	resp, err := svc.Scan(params)

	if err != nil {
		log.Println(err.Error())
	}

	slc := make([]Topic, len(resp.Items))
	lastId = 0
	for key, val := range resp.Items {
		err := dynamodbattribute.UnmarshalMap(val, &slc[key])
		if err != nil {
			log.Println("Error Unmarshaling.", err)
		}
		if slc[key].ID >= lastId {
			lastId = slc[key].ID
		}
	}

	return slc
}
func putItem(w http.ResponseWriter, r *http.Request) {
	id := getLastItemID()+1
	log.Println(reflect.TypeOf(getLastItemID))
	log.Println(id)
	t := Topic{ID: id, Title: r.FormValue("title"),
    Description: r.FormValue("description"),
    Likes: 0, Comments: []Comment{}, Created: time.Now()}

	av, err := dynamodbattribute.Marshal(t)
	log.Println(av, err)

	sess := session.New(&aws.Config{Region: aws.String("us-west-2")})
	if err != nil {
	    log.Println("Failed create session", err)
	    return
	}

	svc := dynamodb.New(sess)
	item, err := dynamodbattribute.MarshalMap(t)
	if err != nil {
	    log.Println("Failed to convert", err)
	    return
	}


	result, err := svc.PutItem(&dynamodb.PutItemInput{
	    Item:      item,
	    TableName: aws.String("discussions"),
	})
	if err != nil {
		log.Println("Failed to save to DynamoDB Table.", err)
		return
	}
	log.Println("Successfully added item to DynamoDB table.")
	log.Println(result)
}

func putComment(r *http.Request) {
	item := getItem(r.FormValue("topic_ID"))
	usr := getName()
	d := ID{ID: item.ID}
	av, err := dynamodbattribute.MarshalMap(d)
	sess := session.New(&aws.Config{Region: aws.String("us-west-2")})

	svc := dynamodb.New(sess)

	c := make([]Comment, len(item.Comments)+1)
	
	if len(item.Comments) > 0 {
		copy(c, item.Comments)
	}

	c[len(c)-1] = Comment{Author: usr, Comment: r.FormValue("comment"), Date: time.Now()}	

	avc, _ := dynamodbattribute.MarshalList(c)
	
	
	params := &dynamodb.UpdateItemInput{TableName: aws.String("discussions"), Key: av, ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{":currComment": {L: avc}}, UpdateExpression: aws.String("SET Comments = :currComment"),}

	resp, err := svc.UpdateItem(params)
	if err != nil {
		log.Println("Error")
		log.Println(err)
	}
	log.Println(resp)
}

func upVoteModel(r *http.Request) {
	item := getItem(r.FormValue("topicID"))

	d := ID{ID: item.ID}
	av, err := dynamodbattribute.MarshalMap(d)
	sess := session.New(&aws.Config{Region: aws.String("us-west-2")})

	svc := dynamodb.New(sess)
	i := 1
	a, _ := dynamodbattribute.Marshal(i)
	params := &dynamodb.UpdateItemInput{TableName: aws.String("discussions"), Key: av, ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{":val": a}, UpdateExpression: aws.String("SET Likes = Likes + :val"),}

	resp, err := svc.UpdateItem(params)
	if err != nil {
		log.Println("Error")
		log.Println(err)
	}
	log.Println(resp)
}

func downVoteModel(r *http.Request) {
	item := getItem(r.FormValue("topicID"))

	d := ID{ID: item.ID}
	av, err := dynamodbattribute.MarshalMap(d)
	sess := session.New(&aws.Config{Region: aws.String("us-west-2")})

	svc := dynamodb.New(sess)
	i := 1
	a, _ := dynamodbattribute.Marshal(i)
	params := &dynamodb.UpdateItemInput{TableName: aws.String("discussions"), Key: av, ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{":val": a}, UpdateExpression: aws.String("SET Likes = Likes - :val"),}

	resp, err := svc.UpdateItem(params)
	if err != nil {
		log.Println("Error")
		log.Println(err)
	}
	log.Println(resp)
}