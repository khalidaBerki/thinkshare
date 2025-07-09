package payment

import (
	"os"

	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/checkout/session"
	"github.com/stripe/stripe-go/v78/price"
	"github.com/stripe/stripe-go/v78/product"
)

func InitStripe() {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
}

// CreateStripeCheckoutSession crée une session de paiement Stripe et retourne l'URL
func CreateStripeCheckoutSession(amount float64, currency, successURL, cancelURL, customerEmail string, metadata map[string]string) (string, string, error) {
	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency:   stripe.String(currency),
					UnitAmount: stripe.Int64(int64(amount * 100)), // attention en centimes
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("Abonnement créateur ThinkShare"),
					},
				},
				Quantity: stripe.Int64(1),
			},
		},
		Mode:          stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL:    stripe.String(successURL),
		CancelURL:     stripe.String(cancelURL),
		CustomerEmail: stripe.String(customerEmail),
	}
	if metadata != nil {
		params.Metadata = metadata
	}

	s, err := session.New(params)
	if err != nil {
		return "", "", err
	}
	return s.ID, s.URL, nil
}

// CreateStripeSubscriptionSession crée une session Stripe Checkout pour un abonnement mensuel
func CreateStripeSubscriptionSession(amount float64, currency, successURL, cancelURL, customerEmail string, metadata map[string]string) (string, string, error) {
	// 1. Créer un produit Stripe (optionnel, on peut réutiliser le même nom)
	prod, err := product.New(&stripe.ProductParams{
		Name: stripe.String("Abonnement créateur ThinkShare"),
	})
	if err != nil {
		return "", "", err
	}

	// 2. Créer un prix récurrent mensuel
	priceParams := &stripe.PriceParams{
		UnitAmount: stripe.Int64(int64(amount * 100)),
		Currency:   stripe.String(currency),
		Recurring: &stripe.PriceRecurringParams{
			Interval: stripe.String("month"),
		},
		Product: stripe.String(prod.ID),
	}
	pr, err := price.New(priceParams)
	if err != nil {
		return "", "", err
	}

	// 3. Créer la session Checkout en mode subscription
	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(pr.ID),
				Quantity: stripe.Int64(1),
			},
		},
		Mode:          stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		SuccessURL:    stripe.String(successURL),
		CancelURL:     stripe.String(cancelURL),
		CustomerEmail: stripe.String(customerEmail),
	}
	if metadata != nil {
		params.Metadata = metadata
	}

	s, err := session.New(params)
	if err != nil {
		return "", "", err
	}
	return s.ID, s.URL, nil
}
