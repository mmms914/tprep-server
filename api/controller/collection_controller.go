package controller

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"main/domain"
	"net/http"
	"slices"
	"strconv"
	"strings"
)

type CollectionController struct {
	CollectionUseCase domain.CollectionUseCase
	UserUseCase       domain.UserUseCase
	HistoryUseCase    domain.HistoryUseCase
}

func (cc *CollectionController) Create(w http.ResponseWriter, r *http.Request) {
	var collection domain.Collection

	err := json.NewDecoder(r.Body).Decode(&collection)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusBadRequest)
		return
	}

	if collection.Name == "" {
		http.Error(w, jsonError("Invalid name"), http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("x-user-id").(string)
	collection.Author = userID

	id, err := cc.CollectionUseCase.Create(r.Context(), &collection, userID)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}

	collection, err = cc.CollectionUseCase.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}

	collectionInfo := domain.CollectionInfo{
		ID:        collection.ID,
		Name:      collection.Name,
		IsPublic:  collection.IsPublic,
		Cards:     collection.Cards,
		Author:    collection.Author,
		Likes:     collection.Likes,
		Trainings: collection.Trainings,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(collectionInfo)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}
}

func (cc *CollectionController) Update(w http.ResponseWriter, r *http.Request) {
	var collection domain.Collection
	userID := r.Context().Value("x-user-id").(string)

	err := json.NewDecoder(r.Body).Decode(&collection)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusBadRequest)
		return
	}

	if collection.Name == "" {
		http.Error(w, jsonError("Invalid name"), http.StatusBadRequest)
		return
	}

	id := chi.URLParam(r, "id")
	coll, err := cc.CollectionUseCase.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, jsonError("There is no collection with this ID"), http.StatusNotFound)
		return
	}

	if coll.Author != userID {
		http.Error(w, jsonError("You are not the owner of this collection"), http.StatusForbidden)
		return
	}

	err = cc.CollectionUseCase.PutByID(r.Context(), id, &collection)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(domain.SuccessResponse{
		Message: "Collection updated",
	})
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}
}

func (cc *CollectionController) AddLike(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID := r.Context().Value("x-user-id").(string)

	user, err := cc.UserUseCase.GetByID(r.Context(), userID)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusNotFound)
	}
	for _, favID := range user.Favourite {
		if favID == id {
			http.Error(w, jsonError("Collection already in favourites"), http.StatusBadRequest)
			return
		}
	}

	coll, err := cc.CollectionUseCase.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, jsonError("There is no collection with this ID"), http.StatusNotFound)
		return
	}

	if !coll.IsPublic && userID != coll.Author {
		http.Error(w, jsonError("Collection is not public"), http.StatusForbidden)
		return
	}

	collection, err := cc.CollectionUseCase.AddLike(r.Context(), id, userID)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]int{"likes": collection.Likes}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}
}

func (cc *CollectionController) RemoveLike(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID := r.Context().Value("x-user-id").(string)

	user, err := cc.UserUseCase.GetByID(r.Context(), userID)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusNotFound)
	}
	found := false
	for _, favID := range user.Favourite {
		if favID == id {
			found = true
			break
		}
	}

	if !found {
		http.Error(w, jsonError("Collection not in favourites"), http.StatusBadRequest)
		return
	}

	coll, err := cc.CollectionUseCase.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, jsonError("There is no collection with this ID"), http.StatusNotFound)
		return
	}

	if !coll.IsPublic && userID != coll.Author {
		http.Error(w, jsonError("Collection is not public"), http.StatusForbidden)
		return
	}

	collection, err := cc.CollectionUseCase.RemoveLike(r.Context(), id, userID)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}

	collectionInfo := domain.CollectionInfo{
		ID:        collection.ID,
		Name:      collection.Name,
		IsPublic:  collection.IsPublic,
		Cards:     collection.Cards,
		Author:    collection.Author,
		Likes:     collection.Likes,
		Trainings: collection.Trainings,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]int{"likes": collectionInfo.Likes}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}
}

func (cc *CollectionController) Get(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("x-user-id").(string)

	id := chi.URLParam(r, "id")
	collection, err := cc.CollectionUseCase.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, jsonError("There is no collection with this ID"), http.StatusNotFound)
		return
	}

	if userID != collection.Author && collection.IsPublic == false {
		http.Error(w, jsonError("You are not the owner of this collection"), http.StatusForbidden)
		return
	}

	for ind, elem := range collection.Cards {
		if elem.OtherAnswers.Items == nil {
			collection.Cards[ind].OtherAnswers.Items = make([]string, 0)
		}
	}

	collectionInfo := domain.CollectionInfo{
		ID:        collection.ID,
		Name:      collection.Name,
		IsPublic:  collection.IsPublic,
		Cards:     collection.Cards,
		Author:    collection.Author,
		Likes:     collection.Likes,
		Trainings: collection.Trainings,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(collectionInfo)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}
}

func (cc *CollectionController) Delete(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("x-user-id").(string)

	id := chi.URLParam(r, "id")

	coll, err := cc.CollectionUseCase.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, jsonError("There is no collection with this ID"), http.StatusNotFound)
		return
	}

	if coll.Author != userID {
		http.Error(w, jsonError("You are not the owner of this collection"), http.StatusForbidden)
		return
	}

	err = cc.CollectionUseCase.DeleteByID(r.Context(), coll.ID, coll.Author)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(domain.SuccessResponse{
		Message: "Collection deleted successfully",
	})
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}
}

func (cc *CollectionController) Search(w http.ResponseWriter, r *http.Request) {
	var collections []domain.Collection
	var result domain.CollectionPreviewArray
	var err error

	userID := r.Context().Value("x-user-id").(string)

	queryParams := r.URL.Query()
	name := queryParams.Get("name")
	count, err := strconv.Atoi(queryParams.Get("count"))
	if err != nil || count < 1 || count > 100 {
		http.Error(w, jsonError("Invalid count"), http.StatusBadRequest)
		return
	}
	offset, err := strconv.Atoi(queryParams.Get("offset"))
	if err != nil || offset < 0 {
		http.Error(w, jsonError("Invalid offset"), http.StatusBadRequest)
		return
	}
	sortBy := queryParams.Get("sort_by")

	if sortBy != "likes" && sortBy != "trainings" && sortBy != "" {
		http.Error(w, jsonError("invalid sort method"), http.StatusBadRequest)
		return
	}
	if sortBy == "" {
		sortBy = "likes"
	}

	category := queryParams.Get("category")
	if category != "" && category != "favourite" {
		http.Error(w, jsonError("invalid category"), http.StatusBadRequest)
		return
	}

	collections, err = cc.CollectionUseCase.SearchPublic(r.Context(), name, count, offset, sortBy, category, userID)

	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}

	if len(collections) == 0 {
		http.Error(w, "Сouldn't find anything", http.StatusNotFound)
		return
	}

	cnt := 0
	for _, c := range collections {
		result.Items = append(result.Items, domain.CollectionPreview{
			ID:         c.ID,
			Name:       c.Name,
			IsPublic:   c.IsPublic,
			CardsCount: len(c.Cards),
			Likes:      c.Likes,
			Trainings:  c.Trainings,
		})
		cnt++
	}
	result.Count = cnt

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}
}

func (cc *CollectionController) CreateCard(w http.ResponseWriter, r *http.Request) {
	var card domain.Card
	userID := r.Context().Value("x-user-id").(string)

	err := json.NewDecoder(r.Body).Decode(&card)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusBadRequest)
		return
	}
	if card.Question == "" || card.Answer == "" || card.OtherAnswers.Items == nil || card.OtherAnswers.Count != len(card.OtherAnswers.Items) {
		http.Error(w, jsonError("Invalid body data"), http.StatusBadRequest)
		return
	}
	collectionID := chi.URLParam(r, "id")

	coll, err := cc.CollectionUseCase.GetByID(r.Context(), collectionID)
	if err != nil {
		http.Error(w, jsonError("There is no collection with this ID"), http.StatusNotFound)
		return
	}

	if coll.Author != userID {
		http.Error(w, jsonError("You are not the owner of this collection"), http.StatusForbidden)
		return
	}

	card, err = cc.CollectionUseCase.AddCard(r.Context(), collectionID, &card)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(card)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}
}

func (cc *CollectionController) UpdateCard(w http.ResponseWriter, r *http.Request) {
	var card domain.Card
	userID := r.Context().Value("x-user-id").(string)

	err := json.NewDecoder(r.Body).Decode(&card)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusBadRequest)
		return
	}

	if card.Question == "" || card.Answer == "" || card.OtherAnswers.Items == nil || card.OtherAnswers.Count != len(card.OtherAnswers.Items) {
		http.Error(w, jsonError("Invalid body data"), http.StatusBadRequest)
		return
	}

	id := chi.URLParam(r, "id")
	card.LocalID, err = strconv.Atoi(chi.URLParam(r, "cardID"))
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusBadRequest)
		return
	}

	coll, err := cc.CollectionUseCase.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, jsonError("There is no collection with this ID"), http.StatusNotFound)
		return
	}

	if coll.Author != userID {
		http.Error(w, jsonError("You are not the owner of this collection"), http.StatusForbidden)
		return
	}

	err = cc.CollectionUseCase.UpdateCard(r.Context(), id, &card)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(domain.SuccessResponse{
		Message: "Card updated",
	})
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}
}

func (cc *CollectionController) DeleteCard(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("x-user-id").(string)

	collectionID := chi.URLParam(r, "id")
	cardID, err := strconv.Atoi(chi.URLParam(r, "cardID"))
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}

	coll, err := cc.CollectionUseCase.GetByID(r.Context(), collectionID)
	if err != nil {
		http.Error(w, jsonError("There is no collection with this ID"), http.StatusNotFound)
		return
	}

	if coll.Author != userID {
		http.Error(w, jsonError("You are not the owner of this collection"), http.StatusForbidden)
		return
	}

	for _, elem := range coll.Cards {
		if elem.LocalID == cardID {
			if elem.Attachment != "" {
				err = cc.CollectionUseCase.RemoveCardPicture(r.Context(), userID, collectionID, cardID, elem.Attachment)
				if err != nil {
					http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
					return
				}
			}
			break
		}
	}

	err = cc.CollectionUseCase.DeleteCard(r.Context(), collectionID, cardID)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(domain.SuccessResponse{
		Message: "Card deleted successfully",
	})
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}
}

func (cc *CollectionController) AddTraining(w http.ResponseWriter, r *http.Request) {
	var historyItem domain.HistoryItem
	userID := r.Context().Value("x-user-id").(string)

	err := json.NewDecoder(r.Body).Decode(&historyItem)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusBadRequest)
		return
	}

	if historyItem.AllCardsCount == 0 || historyItem.CollectionName == "" || historyItem.CollectionID == "" || historyItem.AllCardsCount < len(historyItem.CorrectCards) {
		http.Error(w, jsonError("Invalid body data. Check all_cards_count/collection_name/collection_id"), http.StatusBadRequest)
		return
	}

	if historyItem.IncorrectCards == nil || historyItem.CorrectCards == nil || historyItem.Errors == nil || historyItem.RightAnswers == nil {
		http.Error(w, jsonError("Invalid body data. Check errors/right_answers/(in)correct_cards"), http.StatusBadRequest)
		return
	}

	if historyItem.Time < 0 {
		http.Error(w, jsonError("Invalid body data. Check time"), http.StatusBadRequest)
		return
	}

	err = cc.HistoryUseCase.AddTraining(r.Context(), userID, historyItem)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(domain.SuccessResponse{
		Message: "History item successfully added",
	})
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}
}

func (cc *CollectionController) GetCardPicture(w http.ResponseWriter, r *http.Request) {
	authID := r.Context().Value("x-user-id").(string)

	queryParams := r.URL.Query()
	objectName := queryParams.Get("object_name")

	if objectName == "" || strings.Count(objectName, "_") != 2 {
		http.Error(w, jsonError("Invalid object_name"), http.StatusBadRequest)
		return
	}

	spl := strings.Split(objectName, "_")
	collectionID := spl[0]

	if collectionID != chi.URLParam(r, "id") {
		http.Error(w, jsonError("invalid collection / object_name"), http.StatusBadRequest)
		return
	}

	coll, err := cc.CollectionUseCase.GetByID(r.Context(), collectionID)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusNotFound)
		return
	}

	if coll.IsPublic == false {
		user, err := cc.UserUseCase.GetByID(r.Context(), authID)
		if err != nil {
			http.Error(w, jsonError(err.Error()), http.StatusNotFound)
			return
		}

		if !slices.Contains(user.Collections, collectionID) {
			http.Error(w, jsonError("Access denied"), http.StatusForbidden)
			return
		}
	}

	fileBytes, err := cc.CollectionUseCase.GetCardPhoto(r.Context(), objectName)
	if err != nil {
		http.Error(w, jsonError("Card picture not found"), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(fileBytes)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}
}

func (cc *CollectionController) UploadCardPicture(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(5 << 20)
	if err != nil {
		http.Error(w, jsonError("Error with file or its max size"), http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if handler.Size > (5 << 20) {
		http.Error(w, jsonError("Image's size should be less than 5 MB"), http.StatusBadRequest)
		return
	}

	if !(strings.HasSuffix(handler.Filename, ".jpg") || strings.HasSuffix(handler.Filename, ".jpeg")) {
		http.Error(w, jsonError("Image's extension should be JPG/JPEG"), http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("x-user-id").(string)
	id := chi.URLParam(r, "id")
	cardID, err := strconv.Atoi(chi.URLParam(r, "cardID"))
	if err != nil {
		http.Error(w, jsonError("Invalid path cardID"), http.StatusInternalServerError)
		return
	}

	// check if attachment already is
	coll, err := cc.CollectionUseCase.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusNotFound)
		return
	}

	for _, elem := range coll.Cards {
		if elem.LocalID == cardID {
			if elem.Attachment != "" {
				err = cc.CollectionUseCase.RemoveCardPicture(r.Context(), userID, id, cardID, elem.Attachment)
				if err != nil {
					http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
					return
				}
			}
			break
		}
	}
	//

	objectName, err := cc.CollectionUseCase.UploadCardPhoto(r.Context(), userID, id, cardID, file, handler.Size)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(domain.UploadCardPhotoResult{
		ObjectName: objectName,
	})
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}
}

func (cc *CollectionController) RemoveCardPicture(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("x-user-id").(string)
	queryParams := r.URL.Query()
	objectName := queryParams.Get("object_name")

	collectionID := chi.URLParam(r, "id")
	cardID, err := strconv.Atoi(chi.URLParam(r, "cardID"))
	if err != nil {
		http.Error(w, jsonError("Invalid path cardID"), http.StatusInternalServerError)
		return
	}

	if objectName == "" || strings.Count(objectName, "_") != 2 {
		http.Error(w, jsonError("Invalid object_name"), http.StatusBadRequest)
		return
	}

	spl := strings.Split(objectName, "_")

	if spl[0] != chi.URLParam(r, "id") {
		http.Error(w, jsonError("invalid collection / object_name"), http.StatusBadRequest)
		return
	}

	user, err := cc.UserUseCase.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusNotFound)
		return
	}

	if !slices.Contains(user.Collections, collectionID) {
		http.Error(w, jsonError("Access denied"), http.StatusBadRequest)
		return
	}

	err = cc.CollectionUseCase.RemoveCardPicture(r.Context(), id, collectionID, cardID, objectName)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(domain.SuccessResponse{
		Message: "Card picture deleted",
	})
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}
}
